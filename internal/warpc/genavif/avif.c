#include <stdlib.h>
#include <string.h>
#include <stdbool.h>
#include <stdio.h>
#include "avif/avif.h"
#include "../deps/parson/parson.h"

void handle_commands(FILE *stream);

int main()
{
    // This will read commands from stdin and write responses to stdout
    // and return 0 when stdin is closed.
    // Any errors gets reported in the RPC response messages.
    handle_commands(stdin);
    return 0;
}

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
    int depth; // Bit depth per channel (8, 10, 12, 16).
    int loopCount;
    int frameCount;
    int *frameDurations;

    // CICP color properties.
    int colorPrimaries;
    int transferCharacteristics;
    int matrixCoefficients;

    // HDR CLLI box (max content/picture-average light level, in cd/m^2).
    int maxCLL;
    int maxPALL;

} InputParams;

typedef struct
{
    float quality;        // between 1 and 100.
    char compression[32]; // "lossy" or "lossless"
    int encoderSpeed;     // 1 (slowest, best) to 10 (fastest). 0 means use default.
    char hint[64];        // drawing, icon, photo, picture, or text. Selects chroma subsampling.
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
            msg.data.params.depth = (int)json_object_get_number(params_object, "depth");
            msg.data.params.loopCount = (int)json_object_get_number(params_object, "loopCount");
            msg.data.params.colorPrimaries = (int)json_object_get_number(params_object, "colorPrimaries");
            msg.data.params.transferCharacteristics = (int)json_object_get_number(params_object, "transferCharacteristics");
            msg.data.params.matrixCoefficients = (int)json_object_get_number(params_object, "matrixCoefficients");
            msg.data.params.maxCLL = (int)json_object_get_number(params_object, "maxCLL");
            msg.data.params.maxPALL = (int)json_object_get_number(params_object, "maxPALL");
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
            msg.data.options.quality = (float)json_object_get_number(options_object, "quality");
            msg.data.options.encoderSpeed = (int)json_object_get_number(options_object, "encoderSpeed");
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
        }
    }

    json_value_free(root_value);
    return msg;
}

static void write_blob(uint32_t id, const uint8_t *data, uint32_t size)
{
    uint8_t output_blob_header[16];
    uint32_t output_blob_id = id;
    uint32_t output_blob_size = size;
    // See https://github.com/bep/textandbinarywriter
    const char magic[] = {'T', 'A', 'K', '3', '5', 'E', 'M', '1'};
    memcpy(output_blob_header, magic, 8);
    memcpy(&output_blob_header[8], &output_blob_id, sizeof(output_blob_id));
    memcpy(&output_blob_header[12], &output_blob_size, sizeof(output_blob_size));

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
        json_object_set_number(params_object, "depth", msg->data.params.depth);
        json_object_set_number(params_object, "colorPrimaries", msg->data.params.colorPrimaries);
        json_object_set_number(params_object, "transferCharacteristics", msg->data.params.transferCharacteristics);
        json_object_set_number(params_object, "matrixCoefficients", msg->data.params.matrixCoefficients);
        json_object_set_number(params_object, "maxCLL", msg->data.params.maxCLL);
        json_object_set_number(params_object, "maxPALL", msg->data.params.maxPALL);
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

// drain_bytes discards n bytes from stream. Used to keep the protocol aligned
// after an error that prevents the blob from being consumed normally.
static void drain_bytes(FILE *stream, size_t n)
{
    uint8_t buf[4096];
    while (n > 0)
    {
        size_t want = n < sizeof(buf) ? n : sizeof(buf);
        size_t got = fread(buf, 1, want, stream);
        if (got == 0)
        {
            break;
        }
        n -= got;
    }
}

// avifFormatForHint maps a content hint to a chroma subsampling format.
// Photographic content tolerates 4:2:0, which roughly halves the encoder's
// memory footprint and the output size. Sharp-edged content (text, icons, line
// art) keeps full 4:4:4 chroma. See issue 14987.
static avifPixelFormat avifFormatForHint(const char *hint)
{
    if (strcmp(hint, "drawing") == 0 || strcmp(hint, "icon") == 0 || strcmp(hint, "text") == 0)
    {
        return AVIF_PIXEL_FORMAT_YUV444;
    }
    return AVIF_PIXEL_FORMAT_YUV420; // photo, picture, and the default.
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
            // Out of memory. Drain the blob from the input stream so the next
            // command stays aligned, then report the error to the client instead
            // of leaving the stream corrupted with no response.
            drain_bytes(stream, (size_t)blob_size);
            OutputMessage err_output = {0};
            err_output.header = input.header;
            snprintf(err_output.header.err, sizeof(err_output.header.err),
                     "out of memory allocating %u bytes for blob data", blob_size);
            write_output_message(&err_output);
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
            avifDecoder *decoder = avifDecoderCreate();
            if (decoder == NULL)
            {
                snprintf(output.header.err, sizeof(output.header.err), "Failed to create AVIF decoder");
                write_output_message(&output);
                goto cleanup;
            }

            decoder->ignoreExif = AVIF_TRUE;
            decoder->ignoreXMP = AVIF_TRUE;
            decoder->maxThreads = 1;
            // dav1d decodes ~2-4x faster than aom under WASM.
            decoder->codecChoice = AVIF_CODEC_CHOICE_DAV1D;
            // Request gain map pixels when present (e.g. Lightroom HDR exports).
            decoder->imageContentToDecode = AVIF_IMAGE_CONTENT_ALL;

            avifResult result = avifDecoderSetIOMemory(decoder, blob_data, blob_size);
            if (result != AVIF_RESULT_OK)
            {
                snprintf(output.header.err, sizeof(output.header.err), "Failed to set IO memory: %s", avifResultToString(result));
                avifDecoderDestroy(decoder);
                write_output_message(&output);
                goto cleanup;
            }

            result = avifDecoderParse(decoder);
            if (result != AVIF_RESULT_OK)
            {
                snprintf(output.header.err, sizeof(output.header.err), "Failed to parse AVIF: %s", avifResultToString(result));
                avifDecoderDestroy(decoder);
                write_output_message(&output);
                goto cleanup;
            }

            // If a gain map is present (e.g. Adobe-style SDR+gainmap HDR from Lightroom),
            // bake it into a single true-HDR image in BT.2020/PQ at 10-bit.
            // The downstream pipeline then sees a normal HDR AVIF; SDR clients/displays
            // tone-map automatically.
            avifBool hasGainMap = (decoder->image->gainMap != NULL && decoder->image->gainMap->image != NULL);
            if (hasGainMap)
            {
                // Need a full frame to apply the gain map.
                result = avifDecoderNextImage(decoder);
                if (result != AVIF_RESULT_OK)
                {
                    snprintf(output.header.err, sizeof(output.header.err), "Failed to decode AVIF for gain map: %s", avifResultToString(result));
                    avifDecoderDestroy(decoder);
                    write_output_message(&output);
                    goto cleanup;
                }

                avifRGBImage outRGB;
                memset(&outRGB, 0, sizeof(outRGB));
                avifRGBImageSetDefaults(&outRGB, decoder->image);
                outRGB.format = AVIF_RGB_FORMAT_RGBA;
                outRGB.depth = 16;
                if (avifRGBImageAllocatePixels(&outRGB) != AVIF_RESULT_OK)
                {
                    snprintf(output.header.err, sizeof(output.header.err), "Failed to allocate RGB pixels for gain map apply");
                    avifDecoderDestroy(decoder);
                    write_output_message(&output);
                    goto cleanup;
                }

                // Target the alternate (full-HDR) endpoint: alternateHdrHeadroom = log2(HDR/SDR).
                float hdrHeadroom = 0.0f;
                const avifUnsignedFraction *h = &decoder->image->gainMap->alternateHdrHeadroom;
                if (h->d != 0)
                {
                    hdrHeadroom = (float)h->n / (float)h->d;
                }

                avifContentLightLevelInformationBox outCLLI = {0};
                avifResult applyResult = avifImageApplyGainMap(
                    decoder->image,
                    decoder->image->gainMap,
                    hdrHeadroom,
                    AVIF_COLOR_PRIMARIES_BT2020,
                    AVIF_TRANSFER_CHARACTERISTICS_PQ,
                    &outRGB,
                    &outCLLI,
                    NULL);
                if (applyResult != AVIF_RESULT_OK)
                {
                    snprintf(output.header.err, sizeof(output.header.err), "Failed to apply gain map: %s", avifResultToString(applyResult));
                    avifRGBImageFreePixels(&outRGB);
                    avifDecoderDestroy(decoder);
                    write_output_message(&output);
                    goto cleanup;
                }

                output.data.params.width = outRGB.width;
                output.data.params.height = outRGB.height;
                // Tell the Go side we're 10-bit HDR so it wraps in NRGBA64 and re-encodes as HDR.
                output.data.params.depth = 10;
                output.data.params.stride = outRGB.rowBytes;
                output.data.params.frameCount = 1;
                output.data.params.colorPrimaries = AVIF_COLOR_PRIMARIES_BT2020;
                output.data.params.transferCharacteristics = AVIF_TRANSFER_CHARACTERISTICS_PQ;
                output.data.params.matrixCoefficients = AVIF_MATRIX_COEFFICIENTS_BT2020_NCL;
                output.data.params.maxCLL = outCLLI.maxCLL;
                output.data.params.maxPALL = outCLLI.maxPALL;

                size_t blob_out_size = (size_t)outRGB.rowBytes * outRGB.height;
                write_output_message(&output);
                write_blob(output.header.id, outRGB.pixels, blob_out_size);

                avifRGBImageFreePixels(&outRGB);
                avifDecoderDestroy(decoder);
                goto cleanup;
            }

            output.data.params.width = decoder->image->width;
            output.data.params.height = decoder->image->height;
            output.data.params.loopCount = decoder->repetitionCount;
            output.data.params.frameCount = decoder->imageCount;
            output.data.params.depth = decoder->image->depth;
            output.data.params.colorPrimaries = decoder->image->colorPrimaries;
            output.data.params.transferCharacteristics = decoder->image->transferCharacteristics;
            output.data.params.matrixCoefficients = decoder->image->matrixCoefficients;

            // For 10+ bit images (HDR), use 16-bit RGB to preserve quality.
            // 8-bit depth: 4 bytes per pixel (RGBA, 1 byte each)
            // 16-bit depth: 8 bytes per pixel (RGBA, 2 bytes each)
            int rgb_depth = (decoder->image->depth > 8) ? 16 : 8;
            uint32_t bytes_per_pixel = (rgb_depth == 16) ? 8 : 4;
            uint32_t row_bytes = decoder->image->width * bytes_per_pixel;

            output.data.params.stride = row_bytes;

            size_t frame_size = (size_t)row_bytes * decoder->image->height;
            size_t all_frames_size = frame_size * decoder->imageCount;
            uint8_t *all_frames_data = malloc(all_frames_size);
            if (all_frames_data == NULL)
            {
                snprintf(output.header.err, sizeof(output.header.err), "Failed to allocate memory for decoded image");
                avifDecoderDestroy(decoder);
                write_output_message(&output);
                goto cleanup;
            }

            if (decoder->imageCount > 1)
            {
                output.data.params.frameDurations = malloc(sizeof(int) * decoder->imageCount);
                if (output.data.params.frameDurations == NULL)
                {
                    snprintf(output.header.err, sizeof(output.header.err), "Failed to allocate memory for frame durations");
                    free(all_frames_data);
                    avifDecoderDestroy(decoder);
                    write_output_message(&output);
                    goto cleanup;
                }
            }

            // Set up RGB conversion parameters.
            avifRGBImage rgb;
            memset(&rgb, 0, sizeof(rgb));
            avifRGBImageSetDefaults(&rgb, decoder->image);
            rgb.format = AVIF_RGB_FORMAT_RGBA;
            rgb.depth = rgb_depth;
            rgb.rowBytes = row_bytes;

            int frame_index = 0;
            avifResult next_result;
            while ((next_result = avifDecoderNextImage(decoder)) == AVIF_RESULT_OK)
            {
                rgb.pixels = all_frames_data + (frame_index * frame_size);

                avifResult conv_result = avifImageYUVToRGB(decoder->image, &rgb);
                if (conv_result != AVIF_RESULT_OK)
                {
                    snprintf(output.header.err, sizeof(output.header.err), "Failed to convert to RGBA: %s", avifResultToString(conv_result));
                    free(all_frames_data);
                    if (output.data.params.frameDurations != NULL)
                    {
                        free(output.data.params.frameDurations);
                        output.data.params.frameDurations = NULL;
                    }
                    avifDecoderDestroy(decoder);
                    write_output_message(&output);
                    goto cleanup;
                }

                if (decoder->imageCount > 1)
                {
                    uint64_t duration_ms = (uint64_t)(decoder->imageTiming.duration * 1000.0);
                    output.data.params.frameDurations[frame_index] = (int)duration_ms;
                }

                frame_index++;
            }

            if (next_result != AVIF_RESULT_NO_IMAGES_REMAINING)
            {
                snprintf(output.header.err, sizeof(output.header.err), "Failed to decode AVIF frame %d: %s", frame_index, avifResultToString(next_result));
                free(all_frames_data);
                if (output.data.params.frameDurations != NULL)
                {
                    free(output.data.params.frameDurations);
                    output.data.params.frameDurations = NULL;
                }
                avifDecoderDestroy(decoder);
                write_output_message(&output);
                goto cleanup;
            }

            if (frame_index != decoder->imageCount)
            {
                snprintf(output.header.err, sizeof(output.header.err), "Decoded %d frames, expected %d", frame_index, decoder->imageCount);
                free(all_frames_data);
                if (output.data.params.frameDurations != NULL)
                {
                    free(output.data.params.frameDurations);
                    output.data.params.frameDurations = NULL;
                }
                avifDecoderDestroy(decoder);
                write_output_message(&output);
                goto cleanup;
            }

            write_output_message(&output);
            write_blob(output.header.id, all_frames_data, all_frames_size);

            free(all_frames_data);
            if (output.data.params.frameDurations != NULL)
            {
                free(output.data.params.frameDurations);
                output.data.params.frameDurations = NULL;
            }
            avifDecoderDestroy(decoder);

            goto cleanup;
        }
        else if (strcmp(input.header.command, "config") == 0)
        {
            avifDecoder *decoder = avifDecoderCreate();
            if (decoder == NULL)
            {
                snprintf(output.header.err, sizeof(output.header.err), "Failed to create AVIF decoder");
                write_output_message(&output);
                goto cleanup;
            }

            decoder->ignoreExif = AVIF_TRUE;
            decoder->ignoreXMP = AVIF_TRUE;
            decoder->maxThreads = 1;
            decoder->codecChoice = AVIF_CODEC_CHOICE_DAV1D;

            avifResult result = avifDecoderSetIOMemory(decoder, blob_data, blob_size);
            if (result != AVIF_RESULT_OK)
            {
                snprintf(output.header.err, sizeof(output.header.err), "Failed to set IO memory: %s", avifResultToString(result));
                avifDecoderDestroy(decoder);
                write_output_message(&output);
                goto cleanup;
            }

            result = avifDecoderParse(decoder);
            if (result != AVIF_RESULT_OK)
            {
                snprintf(output.header.err, sizeof(output.header.err), "Failed to parse AVIF: %s", avifResultToString(result));
                avifDecoderDestroy(decoder);
                write_output_message(&output);
                goto cleanup;
            }

            output.data.params.width = decoder->image->width;
            output.data.params.height = decoder->image->height;
            output.data.params.depth = decoder->image->depth;
            output.data.params.loopCount = decoder->repetitionCount;
            output.data.params.frameCount = decoder->imageCount;

            avifDecoderDestroy(decoder);

            write_output_message(&output);
        }
        else if (strcmp(input.header.command, "encodeNRGBA") == 0)
        {
            int width = input.data.params.width;
            int height = input.data.params.height;
            int depth = input.data.params.depth;
            int stride = input.data.params.stride;

            // Default depth to 8 if not specified.
            if (depth == 0) {
                depth = 8;
            }
            // For HDR (10+ bit), input comes as 16-bit RGBA.
            int rgb_depth = (depth > 8) ? 16 : 8;
            int bytes_per_pixel = (rgb_depth > 8) ? 8 : 4;
            if (stride == 0) {
                stride = width * bytes_per_pixel;
            }
            float quality = input.data.options.quality;
            const char* compression = input.data.options.compression;

            if (width == 0 || height == 0) {
                snprintf(output.header.err, sizeof(output.header.err), "encodeNRGBA: width and height must be > 0");
                write_output_message(&output);
                goto cleanup;
            }

            // Pick chroma subsampling from the content hint. Lossless keeps 4:4:4,
            // since subsampling discards chroma and would defeat it.
            avifPixelFormat yuvFormat = avifFormatForHint(input.data.options.hint);
            if (strcmp(compression, "lossless") == 0) {
                yuvFormat = AVIF_PIXEL_FORMAT_YUV444;
            }

            // Create image with the target bit depth for encoding.
            avifImage *image = avifImageCreate(width, height, depth, yuvFormat);
            if (!image) {
                snprintf(output.header.err, sizeof(output.header.err), "encodeNRGBA: Failed to create avifImage");
                write_output_message(&output);
                goto cleanup;
            }

            // Set color properties from input if provided, otherwise use depth-based defaults.
            // This preserves the original color space when re-encoding (e.g., BT.709 SDR vs BT.2020 HDR).
            if (input.data.params.colorPrimaries > 0) {
                image->colorPrimaries = input.data.params.colorPrimaries;
                image->transferCharacteristics = input.data.params.transferCharacteristics;
                image->matrixCoefficients = input.data.params.matrixCoefficients;
            } else if (depth > 8) {
                // Default for 10-bit+: BT.2020/PQ (HDR).
                image->colorPrimaries = AVIF_COLOR_PRIMARIES_BT2020;
                image->transferCharacteristics = AVIF_TRANSFER_CHARACTERISTICS_PQ;
                image->matrixCoefficients = AVIF_MATRIX_COEFFICIENTS_BT2020_NCL;
            } else {
                // Default for 8-bit: BT.709/sRGB (SDR).
                image->colorPrimaries = AVIF_COLOR_PRIMARIES_BT709;
                image->transferCharacteristics = AVIF_TRANSFER_CHARACTERISTICS_SRGB;
                image->matrixCoefficients = AVIF_MATRIX_COEFFICIENTS_BT601;
            }
            image->yuvRange = AVIF_RANGE_FULL;

            // CLLI carries HDR peak/avg light levels from a baked gain map.
            if (input.data.params.maxCLL > 0 || input.data.params.maxPALL > 0) {
                image->clli.maxCLL = (uint16_t)input.data.params.maxCLL;
                image->clli.maxPALL = (uint16_t)input.data.params.maxPALL;
            }

            avifRGBImage rgb;
            avifRGBImageSetDefaults(&rgb, image);
            rgb.format = AVIF_RGB_FORMAT_RGBA;
            rgb.depth = rgb_depth;
            rgb.pixels = (uint8_t *)blob_data;
            rgb.rowBytes = stride;

            avifResult result = avifImageRGBToYUV(image, &rgb);
            if (result != AVIF_RESULT_OK) {
                snprintf(output.header.err, sizeof(output.header.err), "encodeNRGBA: Failed to convert to YUV: %s", avifResultToString(result));
                avifImageDestroy(image);
                write_output_message(&output);
                goto cleanup;
            }

            avifEncoder *encoder = avifEncoderCreate();
            if (!encoder) {
                 snprintf(output.header.err, sizeof(output.header.err), "encodeNRGBA: Failed to create encoder");
                avifImageDestroy(image);
                write_output_message(&output);
                goto cleanup;
            }
			encoder->codecChoice = AVIF_CODEC_CHOICE_AOM;

            if (strcmp(compression, "lossless") == 0) {
                encoder->quality = AVIF_QUALITY_LOSSLESS;
                encoder->qualityAlpha = AVIF_QUALITY_LOSSLESS;
            } else {
                // Map Hugo quality 1-100 to libavif quality 0-100 (higher is better).
                if (quality < 1) quality = 1;
                if (quality > 100) quality = 100;
                int avif_quality = (int)((quality - 1.0) / 99.0 * 100.0);
                // Keep lossy strictly below lossless so quality 100 stays lossy. See issue 14981.
                if (avif_quality >= AVIF_QUALITY_LOSSLESS) avif_quality = AVIF_QUALITY_LOSSLESS - 1;
                encoder->quality = avif_quality;
                encoder->qualityAlpha = avif_quality;
            }
            // Range 0 (slowest, best quality) to 10 (fastest).
            encoder->speed = (input.data.options.encoderSpeed >= 1 && input.data.options.encoderSpeed <= 10)
                ? input.data.options.encoderSpeed
                : 10;
            encoder->autoTiling = AVIF_TRUE;

            result = avifEncoderAddImage(encoder, image, 1, AVIF_ADD_IMAGE_FLAG_SINGLE);
            if (result != AVIF_RESULT_OK) {
                snprintf(output.header.err, sizeof(output.header.err), "encodeNRGBA: Failed to add image to encoder: %s", avifResultToString(result));
                avifEncoderDestroy(encoder);
                avifImageDestroy(image);
                write_output_message(&output);
                goto cleanup;
            }

            avifRWData raw = { NULL, 0 };
            result = avifEncoderFinish(encoder, &raw);
            if (result != AVIF_RESULT_OK) {
                snprintf(output.header.err, sizeof(output.header.err), "encodeNRGBA: Failed to finish encoding: %s", avifResultToString(result));
                avifEncoderDestroy(encoder);
                avifImageDestroy(image);
                write_output_message(&output);
                goto cleanup;
            }

            write_output_message(&output);
            write_blob(output.header.id, raw.data, raw.size);

            avifRWDataFree(&raw);
            avifEncoderDestroy(encoder);
            avifImageDestroy(image);
            
            goto cleanup;
        }
        else if (strcmp(input.header.command, "encodeGray") == 0)
        {
            int width = input.data.params.width;
            int height = input.data.params.height;
            int depth = input.data.params.depth;
            int stride = input.data.params.stride;

            // Default depth to 8 if not specified.
            if (depth == 0) {
                depth = 8;
            }
            int bytes_per_sample = (depth > 8) ? 2 : 1;
            if (stride == 0) {
                stride = width * bytes_per_sample;
            }
            float quality = input.data.options.quality;
            const char* compression = input.data.options.compression;

            if (width == 0 || height == 0) {
                snprintf(output.header.err, sizeof(output.header.err), "encodeGray: width and height must be > 0");
                write_output_message(&output);
                goto cleanup;
            }

            avifImage *image = avifImageCreate(width, height, depth, AVIF_PIXEL_FORMAT_YUV400);
            if (!image) {
                snprintf(output.header.err, sizeof(output.header.err), "encodeGray: Failed to create avifImage");
                write_output_message(&output);
                goto cleanup;
            }

            // Set color properties from input if provided, otherwise use depth-based defaults.
            if (input.data.params.colorPrimaries > 0) {
                image->colorPrimaries = input.data.params.colorPrimaries;
                image->transferCharacteristics = input.data.params.transferCharacteristics;
                image->matrixCoefficients = input.data.params.matrixCoefficients;
            } else if (depth > 8) {
                // Default for 10-bit+: BT.2020/PQ (HDR).
                image->colorPrimaries = AVIF_COLOR_PRIMARIES_BT2020;
                image->transferCharacteristics = AVIF_TRANSFER_CHARACTERISTICS_PQ;
                image->matrixCoefficients = AVIF_MATRIX_COEFFICIENTS_BT2020_NCL;
            } else {
                // Default for 8-bit: BT.709/sRGB (SDR).
                image->colorPrimaries = AVIF_COLOR_PRIMARIES_BT709;
                image->transferCharacteristics = AVIF_TRANSFER_CHARACTERISTICS_SRGB;
                image->matrixCoefficients = AVIF_MATRIX_COEFFICIENTS_BT601;
            }
            image->yuvRange = AVIF_RANGE_FULL;

            if (input.data.params.maxCLL > 0 || input.data.params.maxPALL > 0) {
                image->clli.maxCLL = (uint16_t)input.data.params.maxCLL;
                image->clli.maxPALL = (uint16_t)input.data.params.maxPALL;
            }

            avifResult alloc_result = avifImageAllocatePlanes(image, AVIF_PLANES_YUV);
            if (alloc_result != AVIF_RESULT_OK) {
                snprintf(output.header.err, sizeof(output.header.err), "encodeGray: Failed to allocate planes: %s", avifResultToString(alloc_result));
                avifImageDestroy(image);
                write_output_message(&output);
                goto cleanup;
            }
            uint8_t *src = blob_data;
            uint8_t *dst = image->yuvPlanes[AVIF_CHAN_Y];
            size_t row_bytes = width * bytes_per_sample;
            for (int i = 0; i < height; ++i) {
                memcpy(dst, src, row_bytes);
                src += stride;
                dst += image->yuvRowBytes[AVIF_CHAN_Y];
            }


            avifEncoder *encoder = avifEncoderCreate();
            if (!encoder) {
                 snprintf(output.header.err, sizeof(output.header.err), "encodeGray: Failed to create encoder");
                avifImageDestroy(image);
                write_output_message(&output);
                goto cleanup;
            }
			encoder->codecChoice = AVIF_CODEC_CHOICE_AOM;

            if (strcmp(compression, "lossless") == 0) {
                encoder->quality = AVIF_QUALITY_LOSSLESS;
            } else {
                // Map Hugo quality 1-100 to libavif quality 0-100 (higher is better).
                if (quality < 1) quality = 1;
                if (quality > 100) quality = 100;
                int avif_quality = (int)((quality - 1.0) / 99.0 * 100.0);
                // Keep lossy strictly below lossless so quality 100 stays lossy. See issue 14981.
                if (avif_quality >= AVIF_QUALITY_LOSSLESS) avif_quality = AVIF_QUALITY_LOSSLESS - 1;
                encoder->quality = avif_quality;
            }
            encoder->qualityAlpha = AVIF_QUALITY_LOSSLESS;
            encoder->speed = (input.data.options.encoderSpeed >= 1 && input.data.options.encoderSpeed <= 10)
                ? input.data.options.encoderSpeed
                : 10;
            encoder->autoTiling = AVIF_TRUE;

            avifResult result = avifEncoderAddImage(encoder, image, 1, AVIF_ADD_IMAGE_FLAG_SINGLE);
            if (result != AVIF_RESULT_OK) {
                snprintf(output.header.err, sizeof(output.header.err), "encodeGray: Failed to add image to encoder: %s", avifResultToString(result));
                avifEncoderDestroy(encoder);
                avifImageDestroy(image);
                write_output_message(&output);
                goto cleanup;
            }

            avifRWData raw = { NULL, 0 };
            result = avifEncoderFinish(encoder, &raw);
            if (result != AVIF_RESULT_OK) {
                snprintf(output.header.err, sizeof(output.header.err), "encodeGray: Failed to finish encoding: %s", avifResultToString(result));
                avifEncoderDestroy(encoder);
                avifImageDestroy(image);
                write_output_message(&output);
                goto cleanup;
            }

            write_output_message(&output);
            write_blob(output.header.id, raw.data, raw.size);

            avifRWDataFree(&raw);
            avifEncoderDestroy(encoder);
            avifImageDestroy(image);
            
            goto cleanup;
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
