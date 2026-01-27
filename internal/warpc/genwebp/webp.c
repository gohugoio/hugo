// Copyright 2025 The Hugo Authors. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

#include <stdlib.h>
#include <string.h>
#include <stdbool.h>
#include <stdio.h>
#include <webp/encode.h>
#include <webp/decode.h>
#include "webp/demux.h"
#include <webp/mux.h>
#include "deps/parson/parson.h"

void handle_commands(FILE *stream);

int main()
{
    // This will read commands from stdin and write responses to stdout
    // and return 0 when stdin is closed.
    // Any errors gets reported in the RPC response messages.
    handle_commands(stdin);
    return 0;
}

// Translation table for WebPEncodingError to string.
// This translation table is replicated in all webp bindings on the web. Great design.
static const char *const kErrorMessages[VP8_ENC_ERROR_LAST] = {
    "OK",
    "OUT_OF_MEMORY: Out of memory allocating objects",
    "BITSTREAM_OUT_OF_MEMORY: Out of memory re-allocating byte buffer",
    "NULL_PARAMETER: NULL parameter passed to function",
    "INVALID_CONFIGURATION: configuration is invalid",
    "BAD_DIMENSION: Bad picture dimension. Maximum width and height "
    "allowed is 16383 pixels.",
    "PARTITION0_OVERFLOW: Partition #0 is too big to fit 512k.\n"
    "To reduce the size of this partition, try using less segments "
    "with the -segments option, and eventually reduce the number of "
    "header bits using -partition_limit. More details are available "
    "in the manual (`man cwebp`)",
    "PARTITION_OVERFLOW: Partition is too big to fit 16M",
    "BAD_WRITE: Picture writer returned an I/O error",
    "FILE_TOO_BIG: File would be too big to fit in 4G",
    "USER_ABORT: encoding abort requested by user"};

typedef struct
{
    int version;
    int id;
    char command[256];
    char err[256];
} Header;

typedef struct
{
    int width;
    int height;
    int stride;
    int loopCount;
    int frameCount;
    int *frameDurations;
    bool hasAlpha;

} InputParams;

typedef struct
{
    float quality;        // between 1 and 100.
    char compression[32]; // "lossy" or "lossless"
    char hint[64];        // drawing, icon, photo, picture, or text. Default is photo.
    int preset;           // preset to use; resolved from hint.

    bool useSharpYuv; // use sharp YUV for better quality.
    int method;       // quality/speed trade-off (0=fast, 6=slower-better). Default is 2.

} InputOptions;

typedef struct
{
    InputOptions options;
    InputParams params;

} InputData;

typedef struct
{
    Header header;
    InputData data;

} InputMessage;

typedef struct
{
    Header header;
    InputData data;

} OutputMessage;

#define MAX_LINE_LENGTH 4096

static uint8_t *encodeNRGBA(WebPConfig *config, const uint8_t *rgba, int width, int height, int stride, size_t *output_size)
{
    WebPPicture pic;
    WebPMemoryWriter wrt;
    int ok;
    if (!WebPPictureInit(&pic))
    {
        fprintf(stderr, "WebPPictureInit failed\n");
        return NULL;
    }

    pic.use_argb = 1;
    pic.width = width;
    pic.height = height;
    pic.writer = WebPMemoryWrite;
    pic.custom_ptr = &wrt;
    WebPMemoryWriterInit(&wrt);
    ok = WebPPictureImportRGBA(&pic, rgba, stride);
    if (ok)
    {
        ok = WebPEncode(config, &pic);
        if (!ok)
        {
            fprintf(stderr, "WebPEncode failed: %d (%s)\n", pic.error_code, kErrorMessages[pic.error_code]);
        }
    }
    else
    {
        fprintf(stderr, "WebPPictureImportRGBA failed: %d (%s)\n", pic.error_code, kErrorMessages[pic.error_code]);
    }
    WebPPictureFree(&pic);
    if (!ok)
    {
        WebPMemoryWriterClear(&wrt);
        return NULL;
    }
    *output_size = wrt.size;
    return wrt.mem;
}

static uint8_t *encodeGray(WebPConfig *config, uint8_t *y, int width, int height, int stride, size_t *output_size)
{
    WebPPicture pic;
    WebPMemoryWriter wrt;

    int ok;
    if (!WebPPictureInit(&pic))
    {
        return NULL;
    }

    pic.use_argb = 0;
    pic.width = width;
    pic.height = height;
    pic.y_stride = stride;
    pic.writer = WebPMemoryWrite;
    pic.custom_ptr = &wrt;
    WebPMemoryWriterInit(&wrt);

    const int uvWidth = (int)(((int64_t)width + 1) >> 1);
    const int uvHeight = (int)(((int64_t)height + 1) >> 1);
    const int uvStride = uvWidth;
    const int uvSize = uvStride * uvHeight;
    const int gray = 128;
    uint8_t *chroma;

    chroma = malloc(uvSize);
    if (!chroma)
    {
        return 0;
    }
    memset(chroma, gray, uvSize);

    pic.y = y;
    pic.u = chroma;
    pic.v = chroma;
    pic.uv_stride = uvStride;

    ok = WebPEncode(config, &pic);

    free(chroma);

    WebPPictureFree(&pic);
    if (!ok)
    {
        WebPMemoryWriterClear(&wrt);
        return NULL;
    }
    *output_size = wrt.size;
    return wrt.mem;
}

static uint8_t *encodeNRGBAAnimated(WebPConfig *config, InputParams params, const uint8_t *all_frames_data, size_t *output_size)
{
    WebPAnimEncoderOptions anim_options;
    WebPAnimEncoderOptionsInit(&anim_options);
    anim_options.anim_params.loop_count = params.loopCount;

    WebPAnimEncoder *enc = WebPAnimEncoderNew(params.width, params.height, &anim_options);
    if (enc == NULL)
    {
        fprintf(stderr, "Error creating WebPAnimEncoder\n");
        return NULL;
    }

    int timestamp = 0;
    size_t frame_rgba_size = (size_t)params.stride * params.height;

    for (int i = 0; i < params.frameCount; i++)
    {
        WebPPicture pic;
        if (!WebPPictureInit(&pic))
        {
            fprintf(stderr, "WebPPictureInit failed\n");
            WebPAnimEncoderDelete(enc);
            return NULL;
        }
        pic.use_argb = 1;
        pic.width = params.width;
        pic.height = params.height;

        const uint8_t *frame_rgba = all_frames_data + i * frame_rgba_size;

        if (!WebPPictureImportRGBA(&pic, frame_rgba, params.stride))
        {
            fprintf(stderr, "WebPPictureImportRGBA failed\n");
            WebPPictureFree(&pic);
            WebPAnimEncoderDelete(enc);
            return NULL;
        }

        if (!WebPAnimEncoderAdd(enc, &pic, timestamp, config))
        {
            fprintf(stderr, "WebPAnimEncoderAdd failed\n");
            WebPPictureFree(&pic);
            WebPAnimEncoderDelete(enc);
            return NULL;
        }
        timestamp += params.frameDurations[i];
        WebPPictureFree(&pic);
    }

    if (!WebPAnimEncoderAdd(enc, NULL, timestamp, config))
    {
        fprintf(stderr, "WebPAnimEncoderAdd failed for final frame\n");
        WebPAnimEncoderDelete(enc);
        return NULL;
    }

    WebPData webp_data_out;
    WebPDataInit(&webp_data_out);
    if (!WebPAnimEncoderAssemble(enc, &webp_data_out))
    {
        fprintf(stderr, "WebPAnimEncoderAssemble failed\n");
        WebPAnimEncoderDelete(enc);
        return NULL;
    }
    WebPAnimEncoderDelete(enc);

    *output_size = webp_data_out.size;
    uint8_t *webp_data = malloc(*output_size);
    if (webp_data == NULL)
    {
        fprintf(stderr, "malloc failed for final webp data\n");
        WebPDataClear(&webp_data_out);
        return NULL;
    }
    memcpy(webp_data, webp_data_out.bytes, *output_size);
    WebPDataClear(&webp_data_out);

    return webp_data;
}

static uint8_t initDecoderConfig(WebPDecoderConfig *config, WebPData data)
{
    if (!WebPInitDecoderConfig(config))
    {
        return 0;
    }
    if (WebPGetFeatures(data.bytes, data.size, &config->input) != VP8_STATUS_OK)
    {
        return 0;
    }
    return 1;
}

static uint8_t initEncoderConfig(WebPConfig *config, InputOptions opts)
{
    if (!WebPConfigInit(config))
    {
        return 0;
    }

    if (!WebPConfigPreset(config, opts.preset, opts.quality))
    {
        return 0;
    }

    if (strcmp(opts.compression, "lossless") == 0)
    {
        // Activate the lossless compression mode with the desired efficiency level
        // between 0 (fastest, lowest compression) and 9 (slower, best compression).
        // A good default level is '6', providing a fair tradeoff between compression
        // speed and final compressed size.
        if (!WebPConfigLosslessPreset(config, 6))
        {
            return 1;
        }
    }

    config->use_sharp_yuv = opts.useSharpYuv ? 1 : 0;
    config->method = opts.method;

    return 1;
}

InputMessage parse_input_message(const char *line)
{
    InputMessage msg = {0};

    JSON_Value *root_value = json_parse_string(line);
    if (root_value == NULL)
    {
        fprintf(stderr, "Error parsing JSON line\n");
        return msg;
    }

    if (json_value_get_type(root_value) != JSONObject)
    {
        fprintf(stderr, "Error: Line did not parse to a valid JSON object\n");
        json_value_free(root_value);
        return msg;
    }

    JSON_Object *root_object = json_value_get_object(root_value);

    JSON_Object *header_object = json_object_get_object(root_object, "header");
    if (header_object != NULL)
    {
        msg.header.version = (int)json_object_get_number(header_object, "version");
        msg.header.id = (int)json_object_get_number(header_object, "id");
        const char *command_str = json_object_get_string(header_object, "command");
        if (command_str != NULL)
        {
            strncpy(msg.header.command, command_str, sizeof(msg.header.command) - 1);
            msg.header.command[sizeof(msg.header.command) - 1] = '\0';
        }
        const char *err_str = json_object_get_string(header_object, "err");
        if (err_str != NULL)
        {
            strncpy(msg.header.err, err_str, sizeof(msg.header.err) - 1);
            msg.header.err[sizeof(msg.header.err) - 1] = '\0';
        }
    }

    JSON_Object *data_object = json_object_get_object(root_object, "data");
    if (data_object != NULL)
    {
        JSON_Object *params_object = json_object_get_object(data_object, "params");
        if (params_object != NULL)
        {
            msg.data.params.width = (int)json_object_get_number(params_object, "width");
            msg.data.params.height = (int)json_object_get_number(params_object, "height");
            msg.data.params.stride = (int)json_object_get_number(params_object, "stride");
            msg.data.params.loopCount = (int)json_object_get_number(params_object, "loopCount");
            JSON_Array *durations_array = json_object_get_array(params_object, "frameDurations");
            if (durations_array != NULL)
            {
                size_t count = json_array_get_count(durations_array);
                msg.data.params.frameCount = count;

                if (count > 0)
                {
                    msg.data.params.frameDurations = malloc(sizeof(int) * count);
                    if (msg.data.params.frameDurations != NULL)
                    {
                        for (size_t i = 0; i < count; i++)
                        {
                            msg.data.params.frameDurations[i] = (int)json_array_get_number(durations_array, i);
                        }
                    }
                    else
                    {
                        // Malloc failed.
                        msg.data.params.frameCount = 0;
                    }
                }
            }
        }
        JSON_Object *options_object = json_object_get_object(data_object, "options");
        if (options_object != NULL)
        {
            msg.data.options.quality = (int)json_object_get_number(options_object, "quality");
            const char *compression_str = json_object_get_string(options_object, "compression");
            if (compression_str != NULL)
            {
                strncpy(msg.data.options.compression, compression_str, sizeof(msg.data.options.compression) - 1);
                msg.data.options.compression[sizeof(msg.data.options.compression) - 1] = '\0';
            }
            const char *hint_str = json_object_get_string(options_object, "hint");
            if (hint_str != NULL)
            {
                strncpy(msg.data.options.hint, hint_str, sizeof(msg.data.options.hint) - 1);
                msg.data.options.hint[sizeof(msg.data.options.hint) - 1] = '\0';
            }
            msg.data.options.useSharpYuv = json_object_get_number(options_object, "useSharpYuv") != 0;
            msg.data.options.method = (int)json_object_get_number(options_object, "method");
            if (msg.data.options.method < 0 || msg.data.options.method > 6)
            {
                msg.data.options.method = 4; // default
            }
            if (msg.data.options.quality < 0 || msg.data.options.quality > 100)
            {
                msg.data.options.quality = 75; // default
            }
            if (strlen(msg.data.options.hint) == 0)
            {
                strncpy(msg.data.options.hint, "photo", sizeof(msg.data.options.hint) - 1);
            }

            // Resolve preset from hint
            if (strcmp(msg.data.options.hint, "picture") == 0)
            {
                msg.data.options.preset = WEBP_PRESET_PICTURE;
            }
            else if (strcmp(msg.data.options.hint, "photo") == 0)
            {
                msg.data.options.preset = WEBP_PRESET_PHOTO;
            }
            else if (strcmp(msg.data.options.hint, "drawing") == 0)
            {
                msg.data.options.preset = WEBP_PRESET_DRAWING;
            }
            else if (strcmp(msg.data.options.hint, "icon") == 0)
            {
                msg.data.options.preset = WEBP_PRESET_ICON;
            }
            else if (strcmp(msg.data.options.hint, "text") == 0)
            {
                msg.data.options.preset = WEBP_PRESET_TEXT;
            }
            else
            {
                // I'm not sure what webp's default is, but photo seems like a good choice.
                msg.data.options.preset = WEBP_PRESET_PHOTO;
            }
        }
    }

    json_value_free(root_value);
    return msg;
}

static void write_blob(uint32_t id, const uint8_t *data, uint32_t size)
{
    uint8_t output_blob_header[16];
    // See https://github.com/bep/textandbinarywriter
    const char magic[] = {'T', 'A', 'K', '3', '5', 'E', 'M', '1'};
    memcpy(output_blob_header, magic, 8);
    *(uint32_t *)&output_blob_header[8] = id;
    *(uint32_t *)&output_blob_header[12] = size;

    fwrite(output_blob_header, 1, sizeof(output_blob_header), stdout);
    fwrite(data, 1, (size_t)size, stdout);
    fflush(stdout);
}

void write_output_message(const OutputMessage *msg)
{
    JSON_Value *root_value = json_value_init_object();
    JSON_Object *root_object = json_value_get_object(root_value);

    // Header
    JSON_Value *header_value = json_value_init_object();
    JSON_Object *header_object = json_value_get_object(header_value);
    json_object_set_value(root_object, "header", header_value);
    json_object_set_number(header_object, "version", msg->header.version);
    json_object_set_number(header_object, "id", msg->header.id);
    json_object_set_string(header_object, "err", msg->header.err);

    // Data
    if (msg->data.params.width > 0)
    {
        JSON_Value *data_value = json_value_init_object();
        JSON_Object *data_object = json_value_get_object(data_value);
        json_object_set_value(root_object, "data", data_value);

        JSON_Value *params_value = json_value_init_object();
        JSON_Object *params_object = json_value_get_object(params_value);
        json_object_set_value(data_object, "params", params_value);

        json_object_set_number(params_object, "width", msg->data.params.width);
        json_object_set_number(params_object, "height", msg->data.params.height);
        json_object_set_number(params_object, "stride", msg->data.params.stride);
        json_object_set_boolean(params_object, "hasAlpha", msg->data.params.hasAlpha);
        if (msg->data.params.frameDurations != NULL)
        {
            JSON_Value *durations_value = json_value_init_array();
            JSON_Array *durations_array = json_value_get_array(durations_value);
            for (int i = 0; i < msg->data.params.frameCount; i++)
            {
                json_array_append_number(durations_array, msg->data.params.frameDurations[i]);
            }

            json_object_set_value(params_object, "frameDurations", durations_value);
            json_object_set_number(params_object, "loopCount", msg->data.params.loopCount);
        }
    }

    char *serialized_string = json_serialize_to_string(root_value);
    fprintf(stdout, "%s\n", serialized_string);
    fflush(stdout);

    json_free_serialized_string(serialized_string);
    json_value_free(root_value);
}

void handle_commands(FILE *stream)
{

    char line[MAX_LINE_LENGTH];

    while (fgets(line, sizeof(line), stream) != NULL)
    {
        InputMessage input = {0};
        uint8_t *blob_data = NULL;
        uint32_t blob_size = 0;

        // Remove newline character if present
        line[strcspn(line, "\n")] = 0;

        if (strlen(line) == 0)
        {
            continue;
        }

        input = parse_input_message(line);

        // Next in stream is a blob header defined in https://github.com/bep/textandbinaryreader
        // T', 'A', 'K', '3', '5', 'E', 'M', '1' id uint32, size uint32
        uint8_t blob_header[16];
        size_t read_bytes = fread(blob_header, 1, sizeof(blob_header), stream);
        if (read_bytes != sizeof(blob_header))
        {
            fprintf(stderr, "Error reading blob header\n");
            goto cleanup;
        }
        uint32_t blob_id = *(uint32_t *)&blob_header[8];
        blob_size = *(uint32_t *)&blob_header[12];
        blob_data = malloc((size_t)blob_size);
        if (blob_data == NULL)
        {
            fprintf(stderr, "[%d] Error allocating memory for blob data\n", blob_id);
            goto cleanup;
        }
        read_bytes = fread(blob_data, 1, (size_t)blob_size, stream);
        if (read_bytes != (size_t)blob_size)
        {
            fprintf(stderr, "[%d]  Error reading blob data (size: %llu read: %zu) \n", blob_id, (unsigned long long)blob_size, read_bytes);
            goto cleanup;
        }

        OutputMessage output = {0};
        output.header = input.header;

        if (strcmp(input.header.command, "decode") == 0)
        {

            WebPData data;
            data.bytes = blob_data;
            data.size = (size_t)blob_size;

            WebPDecoderConfig config;
            if (!initDecoderConfig(&config, data))
            {
                strncpy(output.header.err, "Failed to initialize WebPDecoderConfig", sizeof(output.header.err) - 1);
                write_output_message(&output);
                goto cleanup;
            }

            output.data.params.width = config.input.width;
            output.data.params.height = config.input.height;
            output.data.params.stride = config.input.width * 4;
            output.data.params.frameDurations = NULL;

            if (config.input.has_animation)
            {
                WebPAnimDecoderOptions dec_options;
                WebPAnimDecoderOptionsInit(&dec_options);
                dec_options.color_mode = MODE_RGBA;

                WebPAnimDecoder *dec = WebPAnimDecoderNew(&data, &dec_options);
                if (dec == NULL)
                {
                    strncpy(output.header.err, "Failed to create WebPAnimDecoder", sizeof(output.header.err) - 1);
                    write_output_message(&output);
                    goto cleanup;
                }

                WebPAnimInfo anim_info;
                if (!WebPAnimDecoderGetInfo(dec, &anim_info))
                {
                    strncpy(output.header.err, "Failed to get animation info", sizeof(output.header.err) - 1);
                    write_output_message(&output);
                    WebPAnimDecoderDelete(dec);
                    goto cleanup;
                }

                output.data.params.width = anim_info.canvas_width;
                output.data.params.height = anim_info.canvas_height;
                output.data.params.stride = anim_info.canvas_width * 4;
                output.data.params.frameCount = anim_info.frame_count;
                output.data.params.loopCount = anim_info.loop_count;
                output.data.params.hasAlpha = true; // Animated WebP always decoded as RGBA

                int *frameDurations = malloc(sizeof(int) * anim_info.frame_count);
                if (frameDurations == NULL)
                {
                    strncpy(output.header.err, "Failed to allocate memory for frame durations", sizeof(output.header.err) - 1);
                    write_output_message(&output);
                    WebPAnimDecoderDelete(dec);
                    goto cleanup;
                }
                output.data.params.frameDurations = frameDurations;

                // Use a demuxer to get frame durations.
                WebPDemuxer *demux = WebPDemux(&data);
                if (demux != NULL)
                {
                    WebPIterator iter;
                    for (uint32_t i = 0; i < anim_info.frame_count; ++i)
                    {
                        if (WebPDemuxGetFrame(demux, i + 1, &iter))
                        {
                            frameDurations[i] = iter.duration;
                            WebPDemuxReleaseIterator(&iter);
                        }
                        else
                        {
                            frameDurations[i] = 0;
                        }
                    }
                    WebPDemuxDelete(demux);
                }

                size_t frame_size = (size_t)anim_info.canvas_width * 4 * anim_info.canvas_height;
                size_t all_frames_size = frame_size * anim_info.frame_count;
                uint8_t *output_buffer = malloc(all_frames_size);

                if (output_buffer == NULL)
                {
                    strncpy(output.header.err, "Failed to allocate memory for frames", sizeof(output.header.err) - 1);
                    write_output_message(&output);
                    WebPAnimDecoderDelete(dec);
                    free(frameDurations);
                    output.data.params.frameDurations = NULL;
                    goto cleanup;
                }

                int frame_index = 0;
                while (WebPAnimDecoderHasMoreFrames(dec))
                {
                    uint8_t *frame_rgba;
                    int timestamp;
                    if (!WebPAnimDecoderGetNext(dec, &frame_rgba, &timestamp))
                    {
                        break;
                    }
                    memcpy(output_buffer + frame_index * frame_size, frame_rgba, frame_size);
                    frame_index++;
                }

                write_output_message(&output);

                write_blob((uint32_t)input.header.id, output_buffer, (uint32_t)all_frames_size);

                WebPAnimDecoderDelete(dec);
                free(output_buffer);
                free(frameDurations);
                output.data.params.frameDurations = NULL;
            }
            else
            {
                uint8_t *output_buffer;
                int bytesPerPixel;

                output.data.params.hasAlpha = config.input.has_alpha;

                if (config.input.has_alpha)
                {
                    output_buffer = WebPDecodeRGBA(blob_data, (size_t)blob_size, &output.data.params.width, &output.data.params.height);
                    bytesPerPixel = 4;
                }
                else
                {
                    output_buffer = WebPDecodeRGB(blob_data, (size_t)blob_size, &output.data.params.width, &output.data.params.height);
                    bytesPerPixel = 3;
                }

                if (output_buffer == NULL)
                {
                    strncpy(output.header.err, "Failed to decode WebP", sizeof(output.header.err) - 1);
                    write_output_message(&output);
                    goto cleanup;
                }

                output.data.params.stride = output.data.params.width * bytesPerPixel;

                write_output_message(&output);

                size_t output_size = (size_t)(output.data.params.stride * output.data.params.height);
                write_blob((uint32_t)input.header.id, output_buffer, (uint32_t)output_size);

                WebPFree(output_buffer);
            }
        }
        else if (strcmp(input.header.command, "config") == 0)
        {
            int width, height;
            if (!WebPGetInfo(blob_data, (size_t)blob_size, &width, &height))
            {
                strncpy(output.header.err, "Failed to get WebP info", sizeof(output.header.err) - 1);
                write_output_message(&output);
                goto cleanup;
            }

            output.data.params.width = width;
            output.data.params.height = height;

            write_output_message(&output);
        }
        else if (strcmp(input.header.command, "encodeNRGBA") == 0)
        {
            WebPConfig config;
            if (!initEncoderConfig(&config, input.data.options))
            {
                strncpy(output.header.err, "Error initializing WebPConfig", sizeof(output.header.err) - 1);
                write_output_message(&output);
                goto cleanup;
            }

            size_t output_size = 0;
            uint8_t *webp_data = NULL;

            if (input.data.params.frameDurations != NULL)
            {
                webp_data = encodeNRGBAAnimated(&config, input.data.params, blob_data, &output_size);
            }
            else
            {
                webp_data = encodeNRGBA(&config, blob_data, input.data.params.width, input.data.params.height, input.data.params.stride, &output_size);
            }

            if (webp_data == NULL)
            {
                strncpy(output.header.err, "Error encoding NRGBA to WebP", sizeof(output.header.err) - 1);
                write_output_message(&output);
                goto cleanup;
            }

            write_output_message(&output);
            write_blob((uint32_t)input.header.id, webp_data, (uint32_t)output_size);
            free(webp_data);
        }
        else if (strcmp(input.header.command, "encodeGray") == 0)
        {
            WebPConfig config;
            if (!initEncoderConfig(&config, input.data.options))
            {
                strncpy(output.header.err, "Error initializing WebPConfig", sizeof(output.header.err) - 1);
                write_output_message(&output);
                goto cleanup;
            }

            size_t output_size = 0;
            uint8_t *webp_data = encodeGray(&config, blob_data, input.data.params.width, input.data.params.height, input.data.params.stride, &output_size);
            if (webp_data == NULL)
            {
                strncpy(output.header.err, "Error encoding Gray to WebP", sizeof(output.header.err) - 1);
                write_output_message(&output);
                goto cleanup;
            }

            write_output_message(&output);
            write_blob((uint32_t)input.header.id, webp_data, (uint32_t)output_size);

            free(webp_data);
        }
        else
        {
            snprintf(output.header.err, sizeof(output.header.err), "Unknown command: %s", input.header.command);
            write_output_message(&output);
        }

    cleanup:
        if (blob_data != NULL)
        {
            free(blob_data);
        }
        if (input.data.params.frameDurations != NULL)
        {
            free(input.data.params.frameDurations);
        }
    }
}
