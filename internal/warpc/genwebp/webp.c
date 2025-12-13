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
#include <stdio.h>
#include <webp/encode.h>
#include <webp/decode.h>
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

} InputOptions;

typedef struct
{
    InputOptions options;

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
        fprintf(stderr, "WebPPictureImportRGBA failed: %d\n", pic.error_code);
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
        JSON_Object *options_object = json_object_get_object(data_object, "options");
        if (options_object != NULL)
        {
            msg.data.options.width = (int)json_object_get_number(options_object, "width");
            msg.data.options.height = (int)json_object_get_number(options_object, "height");
            msg.data.options.stride = (int)json_object_get_number(options_object, "stride");
        }
    }

    json_value_free(root_value);
    return msg;
}

static void write_blob(uint32_t id, const uint8_t *data, uint32_t size)
{
    uint8_t output_blob_header[16];
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
    if (msg->data.options.width > 0)
    {
        fprintf(stdout, "{\"header\":{\"version\":%d,\"id\":%d,\"err\":\"%s\"},\"data\":{\"options\":{\"width\":%d,\"height\":%d,\"stride\":%d}}}\n",
                msg->header.version,
                msg->header.id,
                msg->header.err,
                msg->data.options.width,
                msg->data.options.height,
                msg->data.options.stride);
    }
    else
    {
        fprintf(stdout, "{\"header\":{\"version\":%d,\"id\":%d,\"err\":\"%s\"}}\n",
                msg->header.version,
                msg->header.id,
                msg->header.err);
    }
    fflush(stdout);
}

void handle_commands(FILE *stream)
{

    char line[MAX_LINE_LENGTH];

    while (fgets(line, sizeof(line), stream) != NULL)
    {
        // Remove newline character if present
        line[strcspn(line, "\n")] = 0;

        if (strlen(line) > 0)
        {
            InputMessage input = parse_input_message(line);

            // Next in stream is a blob header defined in https://github.com/bep/textandbinaryreader
            // T', 'A', 'K', '3', '5', 'E', 'M', '1' id uint32, size uint64
            uint8_t blob_header[16];
            size_t read_bytes = fread(blob_header, 1, sizeof(blob_header), stream);
            if (read_bytes != sizeof(blob_header))
            {
                fprintf(stderr, "Error reading blob header\n");
                continue;
            }
            uint32_t blob_id = *(uint32_t *)&blob_header[8];
            uint64_t blob_size = *(uint64_t *)&blob_header[12];
            uint8_t *blob_data = malloc((size_t)blob_size);
            if (blob_data == NULL)
            {
                fprintf(stderr, "Error allocating memory for blob data\n");
                continue;
            }
            read_bytes = fread(blob_data, 1, (size_t)blob_size, stream);
            if (read_bytes != blob_size)
            {
                // TODO1 redo these error handling to use a goto Err.
                fprintf(stderr, "Error reading blob data\n");
                free(blob_data);
                continue;
            }

            OutputMessage output = {0};
            output.header = input.header;

            if (strcmp(input.header.command, "decode") == 0)
            {
                int width, height;
                if (!WebPGetInfo(blob_data, (size_t)blob_size, &width, &height))
                {
                    strncpy(output.header.err, "Failed to get WebP info", sizeof(output.header.err) - 1);
                    write_output_message(&output);
                    free(blob_data);
                    continue;
                }

                uint8_t *output_buffer = WebPDecodeRGBA(blob_data, (size_t)blob_size, &width, &height);

                if (output_buffer == NULL)
                {
                    strncpy(output.header.err, "Failed to decode WebP", sizeof(output.header.err) - 1);
                    write_output_message(&output);
                    free(blob_data);
                    continue;
                }

                output.data.options.width = width;
                output.data.options.height = height;
                output.data.options.stride = width * 4;

                write_output_message(&output);

                size_t output_size = (size_t)width * height * 4;
                write_blob((uint32_t)input.header.id, output_buffer, (uint32_t)output_size);

                WebPFree(output_buffer);
                free(blob_data);
            }
            else if (strcmp(input.header.command, "config") == 0)
            {
                int width, height;
                if (!WebPGetInfo(blob_data, (size_t)blob_size, &width, &height))
                {
                    strncpy(output.header.err, "Failed to get WebP info", sizeof(output.header.err) - 1);
                    write_output_message(&output);
                    free(blob_data);
                    continue;
                }

                output.data.options.width = width;
                output.data.options.height = height;

                write_output_message(&output);
                free(blob_data);
            }
            else if (strcmp(input.header.command, "encodeNRGBA") == 0)
            {
                WebPConfig config;
                if (!WebPConfigInit(&config))
                {
                    strncpy(output.header.err, "Error initializing WebPConfig", sizeof(output.header.err) - 1);
                    write_output_message(&output);
                    free(blob_data);
                    continue;
                }

                config.lossless = 0;
                config.quality = 75.0f;

                size_t output_size = 0;
                uint8_t *webp_data = encodeNRGBA(&config, blob_data, input.data.options.width, input.data.options.height, input.data.options.stride, &output_size);
                if (webp_data == NULL)
                {
                    strncpy(output.header.err, "Error encoding NRGBA to WebP", sizeof(output.header.err) - 1);
                    write_output_message(&output);
                    free(blob_data);
                    continue;
                }

                write_output_message(&output);
                // Write output to stdout with a blob header as defined in https://github.com/bep/textandbinaryreader
                // T', 'A', 'K', '3', '5', 'E', 'M', '1' id uint32, size uint64
                write_blob((uint32_t)input.header.id, webp_data, (uint32_t)output_size);
                free(webp_data);
                free(blob_data);
            }
            else if (strcmp(input.header.command, "encodeGray") == 0)
            {
                WebPConfig config;
                if (!WebPConfigInit(&config))
                {
                    strncpy(output.header.err, "Error initializing WebPConfig", sizeof(output.header.err) - 1);
                    write_output_message(&output);
                    free(blob_data);
                    continue;
                }

                config.lossless = 0;
                config.quality = 75.0f;

                size_t output_size = 0;
                uint8_t *webp_data = encodeGray(&config, blob_data, input.data.options.width, input.data.options.height, input.data.options.stride, &output_size);
                if (webp_data == NULL)
                {
                    strncpy(output.header.err, "Error encoding Gray to WebP", sizeof(output.header.err) - 1);
                    write_output_message(&output);
                    free(blob_data);
                    continue;
                }

                write_output_message(&output);
                // Write output to stdout with a blob header as defined in https://github.com/bep/textandbinaryreader
                // T', 'A', 'K', '3', '5', 'E', 'M', '1' id uint32, size uint64
                write_blob((uint32_t)input.header.id, webp_data, (uint32_t)output_size);

                fprintf(stderr, "Encoded WebP size: %zu bytes\n", output_size);
                free(webp_data);
                free(blob_data);
            }
            else
            {
                snprintf(output.header.err, sizeof(output.header.err), "Unknown command: %s", input.header.command);
                write_output_message(&output);
                free(blob_data);
            }
        }
    }
}
