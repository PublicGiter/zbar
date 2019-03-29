package zbar

// #cgo LDFLAGS: -lzbar
// #cgo CFLAGS: -IE:/SoftWare/ZBar/include
// #cgo 386   LDFLAGS: -L E:/SoftWare/ZBarWin64/lib   -lzbar-0
// #cgo amd64 LDFLAGS: -L E:/SoftWare/ZBarWin64/lib   -lzbar64-0
// #include <zbar.h>
import "C"
import (
	"fmt"
	"reflect"
	"unsafe"
	"syscall"
)

/** "color" of element: bar or space. */
type ZBarColor int
const (
	ZBAR_SPACE   ZBarColor = iota     /**< light area or space between bars */
	ZBAR_BAR                           /**< dark area or colored bar segment */
)

/** decoded symbol type. */
type ZBarSymbolType int
const (
	ZBAR_NONE        ZBarSymbolType =      0 + iota  /**< no symbol decoded */
	ZBAR_PARTIAL     ZBarSymbolType =      1  /**< intermediate status */
	ZBAR_EAN8        ZBarSymbolType =      8  /**< EAN-8 */
	ZBAR_UPCE        ZBarSymbolType =      9  /**< UPC-E */
	ZBAR_ISBN10      ZBarSymbolType =     10  /**< ISBN-10 (from EAN-13). @since 0.4 */
	ZBAR_UPCA        ZBarSymbolType =     12  /**< UPC-A */
	ZBAR_EAN13       ZBarSymbolType =     13  /**< EAN-13 */
	ZBAR_ISBN13      ZBarSymbolType =     14  /**< ISBN-13 (from EAN-13). @since 0.4 */
	ZBAR_I25         ZBarSymbolType =     25  /**< Interleaved 2 of 5. @since 0.4 */
	ZBAR_CODE39      ZBarSymbolType =     39  /**< Code 39. @since 0.4 */
	ZBAR_PDF417      ZBarSymbolType =     57  /**< PDF417. @since 0.6 */
	ZBAR_QRCODE      ZBarSymbolType =     64  /**< QR Code. @since 0.10 */
	ZBAR_CODE128     ZBarSymbolType =    128  /**< Code 128 */
	ZBAR_SYMBOL      ZBarSymbolType = 0x00ff  /**< mask for base symbol type */
	ZBAR_ADDON2      ZBarSymbolType = 0x0200  /**< 2-digit add-on flag */
	ZBAR_ADDON5      ZBarSymbolType = 0x0500  /**< 5-digit add-on flag */
	ZBAR_ADDON       ZBarSymbolType = 0x0700  /**< add-on flag mask */
)

/** error codes. */
type ZBarError int
const (
	ZBAR_OK  ZBarError = iota   /**< no error */
	ZBAR_ERR_NOMEM             /**< out of memory */
	ZBAR_ERR_INTERNAL          /**< internal library error */
	ZBAR_ERR_UNSUPPORTED       /**< unsupported request */
	ZBAR_ERR_INVALID           /**< invalid request */
	ZBAR_ERR_SYSTEM            /**< system error */
	ZBAR_ERR_LOCKING           /**< locking error */
	ZBAR_ERR_BUSY              /**< all resources busy */
	ZBAR_ERR_XDISPLAY          /**< X11 display error */
	ZBAR_ERR_XPROTO            /**< X11 protocol error */
	ZBAR_ERR_CLOSED            /**< output window is closed */
	ZBAR_ERR_WINAPI            /**< windows system error */
	ZBAR_ERR_NUM                /**< number of error codes */
)

/** decoder configuration options.
 * @since 0.4
 */
type ZBarConfig int
const (
	ZBAR_CFG_ENABLE ZBarConfig = iota                /**< enable symbology/feature */
	ZBAR_CFG_ADD_CHECK                               /**< enable check digit when optional */
	ZBAR_CFG_EMIT_CHECK                              /**< return check digit when present */
	ZBAR_CFG_ASCII                                    /**< enable full ASCII character set */
	ZBAR_CFG_NUM                                      /**< number of boolean decoder configs */
)
const (
	ZBAR_CFG_MIN_LEN ZBarConfig = 0x20 + iota        /**< minimum data length for valid decode */
	ZBAR_CFG_MAX_LEN                                 /**< maximum data length for valid decode */
	ZBAR_CFG_POSITION ZBarConfig = 0x80              /**< enable scanner to collect position data */
)
const (
	ZBAR_CFG_X_DENSITY ZBarConfig = 0x100 + iota    /**< image scanner vertical scan density */
	ZBAR_CFG_Y_DENSITY                               /**< image scanner horizontal scan density */
)

/** retrieve runtime library version information.
 * @param major set to the running major version (unless NULL)
 * @param minor set to the running minor version (unless NULL)
 * @returns 0
 */
func ZBarVersion(major, minor *uint32) int {
	return int(C.zbar_version((*C.uint)(major), (*C.uint)(minor)))
}

/** set global library debug level.
 * @param verbosity desired debug level.  higher values create more spew
 */
func ZBarSetVerbosity(verbosity int) {
	C.zbar_set_verbosity(C.int(verbosity))
}

/** increase global library debug level.
 * eg, for -vvvv
 */
func ZBarIncreaseVerbosity() {
	C.zbar_increase_verbosity()
}

/** retrieve string name for symbol encoding.
 * @param sym symbol type encoding
 * @returns the static string name for the specified symbol type,
 * or "UNKNOWN" if the encoding is not recognized
 */
func ZBarGetSymbolName(sym ZBarSymbolType) string {
	return C.GoString(C.zbar_get_symbol_name(C.zbar_symbol_type_t(sym)))
}

/** retrieve string name for addon encoding.
 * @param sym symbol type encoding
 * @returns static string name for any addon, or the empty string
 * if no addons were decoded
 */
func ZBarGetAddonName(sym ZBarSymbolType) string {
	return C.GoString(C.zbar_get_addon_name(C.zbar_symbol_type_t(sym)))
}

/** parse a configuration string of the form "[symbology.]config[=value]".
 * the config must match one of the recognized names.
 * the symbology, if present, must match one of the recognized names.
 * if symbology is unspecified, it will be set to 0.
 * if value is unspecified it will be set to 1.
 * @returns 0 if the config is parsed successfully, 1 otherwise
 * @since 0.4
 */
func ZBarParseConfig(configString string, symbology *ZBarSymbolType, config *ZBarConfig, value *int) int {
	return int(C.zbar_parse_config(C.CString(configString), (*C.zbar_symbol_type_t)(unsafe.Pointer(symbology)), (*C.zbar_config_t)(unsafe.Pointer(config)), (*C.int)(unsafe.Pointer(value))))
}

/** @internal type unsafe error API (don't use) */
func ZBarErrorSpew(object unsafe.Pointer, verbosity int) int {
	return int(C._zbar_error_spew(object, C.int(verbosity)))
}

func ZBarErrorString(object unsafe.Pointer, verbosity int) string {
	return C.GoString(C._zbar_error_string(object, C.int(verbosity)))
}

func ZBarGetErrorCode(object unsafe.Pointer) ZBarError {
	return ZBarError(C._zbar_get_error_code(object))
}

/*@}*/

type ZBarSymbol struct {}
type ZBarSymbolSet struct {}


/*------------------------------------------------------------*/
/** @name Symbol interface
 * decoded barcode symbol result object.  stores type, data, and image
 * location of decoded symbol.  all memory is owned by the library
 */
/*@{*/

/** @typedef zbar_symbol_t
 * opaque decoded symbol object.
 */

/** symbol reference count manipulation.
 * increment the reference count when you store a new reference to the
 * symbol.  decrement when the reference is no longer used.  do not
 * refer to the symbol once the count is decremented and the
 * containing image has been recycled or destroyed.
 * @note the containing image holds a reference to the symbol, so you
 * only need to use this if you keep a symbol after the image has been
 * destroyed or reused.
 * @since 0.9
 */
func ZBarSymbolRef(symbol *ZBarSymbol, refs int) {
	C.zbar_symbol_ref((*C.zbar_symbol_t)(unsafe.Pointer(symbol)), C.int(refs))
}

/** retrieve type of decoded symbol.
 * @returns the symbol type
 */
func ZBarSymbolGetType(symbol *ZBarSymbol) ZBarSymbolType {
	return ZBarSymbolType(C.zbar_symbol_get_type((*C.struct_zbar_symbol_s)(unsafe.Pointer(symbol))))
}

/** retrieve data decoded from symbol.
 * @returns the data string
 */
func ZBarSymbolGetData(symbol *ZBarSymbol) string {
	return C.GoString(C.zbar_symbol_get_data((*C.zbar_symbol_t)(unsafe.Pointer(symbol))))
}

/** retrieve length of binary data.
 * @returns the length of the decoded data
 */
func ZBarSymbolGetDataLength(symbol *ZBarSymbol) uint32 {
	return uint32(C.zbar_symbol_get_data_length((*C.zbar_symbol_t)(unsafe.Pointer(symbol))))
}

/** retrieve a symbol confidence metric.
 * @returns an unscaled, relative quantity: larger values are better
 * than smaller values, where "large" and "small" are application
 * dependent.
 * @note expect the exact definition of this quantity to change as the
 * metric is refined.  currently, only the ordered relationship
 * between two values is defined and will remain stable in the future
 * @since 0.9
 */
func ZBarSymbolGetQuality(symbol *ZBarSymbol) int {
	return int(C.zbar_symbol_get_quality((*C.zbar_symbol_t)(unsafe.Pointer(symbol))))
}

/** retrieve current cache count.  when the cache is enabled for the
 * image_scanner this provides inter-frame reliability and redundancy
 * information for video streams.
 * @returns < 0 if symbol is still uncertain.
 * @returns 0 if symbol is newly verified.
 * @returns > 0 for duplicate symbols
 */
func ZBarSymbolGetCount(symbol *ZBarSymbol) int {
	return int(C.zbar_symbol_get_count((*C.zbar_symbol_t)(unsafe.Pointer(symbol))))
}

/** retrieve the number of points in the location polygon.  the
 * location polygon defines the image area that the symbol was
 * extracted from.
 * @returns the number of points in the location polygon
 * @note this is currently not a polygon, but the scan locations
 * where the symbol was decoded
 */
func ZBarSymbolGetLocSize(symbol *ZBarSymbol) uint32 {
	return uint32(C.zbar_symbol_get_loc_size((*C.zbar_symbol_t)(unsafe.Pointer(symbol))))
}

/** retrieve location polygon x-coordinates.
 * points are specified by 0-based index.
 * @returns the x-coordinate for a point in the location polygon.
 * @returns -1 if index is out of range
 */
func ZBarSymbolGetLocX(symbol *ZBarSymbol, index uint32) int {
	return int(C.zbar_symbol_get_loc_x((*C.zbar_symbol_t)(unsafe.Pointer(symbol)), C.uint(index)))
}

/** retrieve location polygon y-coordinates.
 * points are specified by 0-based index.
 * @returns the y-coordinate for a point in the location polygon.
 * @returns -1 if index is out of range
 */
func ZBarSymbolGetLocY(symbol *ZBarSymbol, index uint32) int {
	return int(C.zbar_symbol_get_loc_y((*C.zbar_symbol_t)(unsafe.Pointer(symbol)), C.uint(index)))
}

/** iterate the set to which this symbol belongs (there can be only one).
 * @returns the next symbol in the set, or
 * @returns NULL when no more results are available
 */
func ZBarSymbolNext(symbol *ZBarSymbol) *ZBarSymbol {
	return (*ZBarSymbol)(C.zbar_symbol_next((*C.zbar_symbol_t)(unsafe.Pointer(symbol))))
}

/** retrieve components of a composite result.
 * @returns the symbol set containing the components
 * @returns NULL if the symbol is already a physical symbol
 * @since 0.10
 */
func ZBarSymbolGetComponents(symbol *ZBarSymbol) *ZBarSymbolSet {
	return (*ZBarSymbolSet)(C.zbar_symbol_get_components((*C.zbar_symbol_t)(unsafe.Pointer(symbol))))
}

/** iterate components of a composite result.
 * @returns the first physical component symbol of a composite result
 * @returns NULL if the symbol is already a physical symbol
 * @since 0.10
 */
func ZBarSymbolFirstComponent(symbol *ZBarSymbol) *ZBarSymbol {
	return (*ZBarSymbol)(C.zbar_symbol_first_component((*C.zbar_symbol_t)(unsafe.Pointer(symbol))))
}

/** print XML symbol element representation to user result buffer.
 * @see http://zbar.sourceforge.net/2008/barcode.xsd for the schema.
 * @param symbol is the symbol to print
 * @param buffer is the inout result pointer, it will be reallocated
 * with a larger size if necessary.
 * @param buflen is inout length of the result buffer.
 * @returns the buffer pointer
 * @since 0.6
 */
func ZBarSymbolXml(symbol *ZBarSymbol, buffer **byte, bufLen *uint32) *byte {
	return (*byte)(unsafe.Pointer(C.zbar_symbol_xml((*C.zbar_symbol_t)(unsafe.Pointer(symbol)), (**C.char)(unsafe.Pointer(buffer)), (*C.uint)(unsafe.Pointer(bufLen)))))
}

/*@}*/

/*------------------------------------------------------------*/
/** @name Symbol Set interface
 * container for decoded result symbols associated with an image
 * or a composite symbol.
 * @since 0.10
 */
/*@{*/

/** @typedef zbar_symbol_set_t
 * opaque symbol iterator object.
 * @since 0.10
 */

/** reference count manipulation.
 * increment the reference count when you store a new reference.
 * decrement when the reference is no longer used.  do not refer to
 * the object any longer once references have been released.
 * @since 0.10
 */
func ZBarSymbolSetRef(symbols *ZBarSymbolSet, refs int) {
	C.zbar_symbol_set_ref((*C.zbar_symbol_set_t)(unsafe.Pointer(symbols)), C.int(refs))
}

/** retrieve set size.
 * @returns the number of symbols in the set.
 * @since 0.10
 */
func ZBarSymbolSetGetSize(symbols *ZBarSymbolSet) int {
	return int(C.zbar_symbol_set_get_size((*C.zbar_symbol_set_t)(unsafe.Pointer(symbols))))
}

/** set iterator.
 * @returns the first decoded symbol result in a set
 * @returns NULL if the set is empty
 * @since 0.10
 */
func ZBarSymbolSetFirstSymbol(symbols *ZBarSymbolSet) *ZBarSymbol {
	return (*ZBarSymbol)(C.zbar_symbol_set_first_symbol((*C.zbar_symbol_set_t)(unsafe.Pointer(symbols))))
}

/*@}*/

/*------------------------------------------------------------*/
/** @name Image interface
 * stores image data samples along with associated format and size
 * metadata
 */
/*@{*/

/** opaque image object. */
type ZBarImage struct{}

/** cleanup handler callback function.
 * called to free sample data when an image is destroyed.
 */
type ZBarImageCleanupHandler func(image *ZBarImage)

/** data handler callback function.
 * called when decoded symbol results are available for an image
 */
type ZBarImageDataHandler func(image *ZBarImage, userData unsafe.Pointer)

/** new image constructor.
 * @returns a new image object with uninitialized data and format.
 * this image should be destroyed (using zbar_image_destroy()) as
 * soon as the application is finished with it
 */
func ZBarImageCreate() *ZBarImage {
	return (*ZBarImage)(C.zbar_image_create())
}

/** image destructor.  all images created by or returned to the
 * application should be destroyed using this function.  when an image
 * is destroyed, the associated data cleanup handler will be invoked
 * if available
 * @note make no assumptions about the image or the data buffer.
 * they may not be destroyed/cleaned immediately if the library
 * is still using them.  if necessary, use the cleanup handler hook
 * to keep track of image data buffers
 */
func ZBarImageDestroy(image *ZBarImage) {
	C.zbar_image_destroy((*C.zbar_image_t)(unsafe.Pointer(image)))
}

/** image reference count manipulation.
 * increment the reference count when you store a new reference to the
 * image.  decrement when the reference is no longer used.  do not
 * refer to the image any longer once the count is decremented.
 * zbar_image_ref(image, -1) is the same as zbar_image_destroy(image)
 * @since 0.5
 */
func ZBarImageRef(image *ZBarImage, refs int) {
	C.zbar_image_ref((*C.zbar_image_t)(unsafe.Pointer(image)), C.int(refs))
}

/** image format conversion.  refer to the documentation for supported
 * image formats
 * @returns a @em new image with the sample data from the original image
 * converted to the requested format.  the original image is
 * unaffected.
 * @note the converted image size may be rounded (up) due to format
 * constraints
 */
func ZBarImageConvert(image *ZBarImage, format uint64) *ZBarImage {
	return (*ZBarImage)(C.zbar_image_convert((*C.zbar_image_t)(unsafe.Pointer(image)), (C.ulong)(format)))
}

/** image format conversion with crop/pad.
 * if the requested size is larger than the image, the last row/column
 * are duplicated to cover the difference.  if the requested size is
 * smaller than the image, the extra rows/columns are dropped from the
 * right/bottom.
 * @returns a @em new image with the sample data from the original
 * image converted to the requested format and size.
 * @note the image is @em not scaled
 * @see zbar_image_convert()
 * @since 0.4
 */
func ZBarImageConvertResize(image *ZBarImage, format uint64, width, height uint32) *ZBarImage {
	return (*ZBarImage)(C.zbar_image_convert_resize((*C.zbar_image_t)(unsafe.Pointer(image)), C.ulong(format), C.uint(width), C.uint(height)))
}

/** retrieve the image format.
 * @returns the fourcc describing the format of the image sample data
 */
func ZBarImageGetFormat(image *ZBarImage) uint64 {
	return uint64(C.zbar_image_get_format((*C.zbar_image_t)(unsafe.Pointer(image))))
}

/** retrieve a "sequence" (page/frame) number associated with this image.
 * @since 0.6
 */
func ZBarImageGetSequence(image *ZBarImage) uint {
	return uint(C.zbar_image_get_sequence((*C.zbar_image_t)(unsafe.Pointer(image))))
}

/** retrieve the width of the image.
 * @returns the width in sample columns
 */
func ZBarImageGetWidth(image *ZBarImage) uint {
	return uint(C.zbar_image_get_width((*C.zbar_image_t)(unsafe.Pointer(image))))
}

/** retrieve the height of the image.
 * @returns the height in sample rows
 */
func ZBarImageGetHeight(image *ZBarImage) uint {
	return uint(C.zbar_image_get_height((*C.zbar_image_t)(unsafe.Pointer(image))))
}

/** return the image sample data.  the returned data buffer is only
 * valid until zbar_image_destroy() is called
 */
func ZBarImageGetData(image *ZBarImage) unsafe.Pointer {
	return C.zbar_image_get_data((*C.zbar_image_t)(unsafe.Pointer(image)))
}

/** return the size of image data.
 * @since 0.6
 */
func ZBarImageGetDataLength(img *ZBarImage) uint64 {
	return uint64(C.zbar_image_get_data_length((*C.zbar_image_t)(unsafe.Pointer(img))))
}

/** retrieve the decoded results.
 * @returns the (possibly empty) set of decoded symbols
 * @returns NULL if the image has not been scanned
 * @since 0.10
 */
func ZBarImageGetSymbols(image *ZBarImage) *ZBarSymbolSet {
	return (*ZBarSymbolSet)(C.zbar_image_get_symbols((*C.zbar_image_t)(unsafe.Pointer(image))))
}

/** associate the specified symbol set with the image, replacing any
 * existing results.  use NULL to release the current results from the
 * image.
 * @see zbar_image_scanner_recycle_image()
 * @since 0.10
 */
func ZBarImageSetSymbols(img *ZBarImage, symbols *ZBarSymbolSet) {
	C.zbar_image_set_symbols((*C.zbar_image_t)(unsafe.Pointer(img)), (*C.zbar_symbol_set_t)(unsafe.Pointer(symbols)))
}

/** image_scanner decode result iterator.
 * @returns the first decoded symbol result for an image
 * or NULL if no results are available
 */
func ZBarImageFirstSymbol(image *ZBarImage) *ZBarSymbol {
	return (*ZBarSymbol)(C.zbar_image_first_symbol((*C.zbar_image_t)(unsafe.Pointer(image))))
}

/** specify the fourcc image format code for image sample data.
 * refer to the documentation for supported formats.
 * @note this does not convert the data!
 * (see zbar_image_convert() for that)
 */
func ZBarImageSetFormat(img *ZBarImage, format uint64) {
	C.zbar_image_set_format((*C.zbar_image_t)(unsafe.Pointer(img)), C.ulong(format))
}

/** associate a "sequence" (page/frame) number with this image.
 * @since 0.6
 */
func ZBarImageSetSequence(img *ZBarImage, sequenceNum uint32) {
	C.zbar_image_set_sequence((*C.zbar_image_t)(unsafe.Pointer(img)), C.uint(sequenceNum))
}

/** specify the pixel size of the image.
 * @note this does not affect the data!
 */
func ZBarImageSetSize(img *ZBarImage, width, height uint32) {
	C.zbar_image_set_size((*C.zbar_image_t)(unsafe.Pointer(img)), C.uint(width), C.uint(height))
}

/** specify image sample data.  when image data is no longer needed by
 * the library the specific data cleanup handler will be called
 * (unless NULL)
 * @note application image data will not be modified by the library
 */
func ZBarImageSetData(image *ZBarImage, data unsafe.Pointer, dataByteLength uint64, cleanupHandler ZBarImageCleanupHandler) {
	// TODO...

	C.zbar_image_set_data((*C.zbar_image_t)(unsafe.Pointer(image)), data, C.ulong(dataByteLength), (*C.zbar_image_cleanup_handler_t)(unsafe.Pointer(syscall.NewCallback(cleanupHandler))))
}

/** built-in cleanup handler.
 * passes the image data buffer to free()
 */
func ZBarImageFreeData(image *ZBarImage) {
	C.zbar_image_free_data((*C.zbar_image_t)(unsafe.Pointer(image)))
}

/** associate user specified data value with an image.
 * @since 0.5
 */
func ZBarImageSetUserData(image *ZBarImage, userData unsafe.Pointer) {
	C.zbar_image_set_userdata((*C.zbar_image_t)(unsafe.Pointer(image)), userData)
}

/** return user specified data value associated with the image.
 * @since 0.5
 */
func ZBarImageGetUserData(image *ZBarImage) unsafe.Pointer {
	return C.zbar_image_get_userdata((*C.zbar_image_t)(unsafe.Pointer(image)))
}

/** dump raw image data to a file for debug.
 * the data will be prefixed with a 16 byte header consisting of:
 *   - 4 bytes uint = 0x676d697a ("zimg")
 *   - 4 bytes format fourcc
 *   - 2 bytes width
 *   - 2 bytes height
 *   - 4 bytes size of following image data in bytes
 * this header can be dumped w/eg:
 * @verbatim
       od -Ax -tx1z -N16 -w4 [file]
@endverbatim
 * for some formats the image can be displayed/converted using
 * ImageMagick, eg:
 * @verbatim
       display -size 640x480+16 [-depth ?] [-sampling-factor ?x?] \
           {GRAY,RGB,UYVY,YUV}:[file]
@endverbatim
 *
 * @param image the image object to dump
 * @param filebase base filename, appended with ".XXXX.zimg" where
 * XXXX is the format fourcc
 * @returns 0 on success or a system error code on failure
 */
func ZBarImageWrite(image *ZBarImage, fileBase string) int {
	return int(C.zbar_image_write((*C.zbar_image_t)(unsafe.Pointer(image)), C.CString(fileBase)))
}

/** read back an image in the format written by zbar_image_write()
 * @note TBD
 */
func ZBarImageRead(filename string) *ZBarImage {
	return (*ZBarImage)(C.zbar_image_read(C.CString(filename)))
}

/*@}*/

/*------------------------------------------------------------*/
/** @name Processor interface
 * @anchor c-processor
 * high-level self-contained image processor.
 * processes video and images for barcodes, optionally displaying
 * images to a library owned output window
 */
/*@{*/

type ZBarProcessorStruct struct{}
/** opaque standalone processor object. */
type ZBarProcessor ZBarProcessorStruct

/** constructor.
 * if threaded is set and threading is available the processor
 * will spawn threads where appropriate to avoid blocking and
 * improve responsiveness
 */
func ZBarProcessorCreate(threaded int) *ZBarProcessor {
	return (*ZBarProcessor)(C.zbar_processor_create(C.int(threaded)))
}

/** destructor.  cleans up all resources associated with the processor
 */
func ZBarProcessorDestroy(processor *ZBarProcessor) {
	C.zbar_processor_destroy((*C.zbar_processor_t)(unsafe.Pointer(processor)))
}

/** (re)initialization.
 * opens a video input device and/or prepares to display output
 */
func ZBarProcessorInit(processor *ZBarProcessor, videoDevice string, enableDisplay int) int {
	return int(C.zbar_processor_init((*C.zbar_processor_t)(unsafe.Pointer(processor)), C.CString(videoDevice), C.int(enableDisplay)))
}

/** request a preferred size for the video image from the device.
 * the request may be adjusted or completely ignored by the driver.
 * @note must be called before zbar_processor_init()
 * @since 0.6
 */
func ZBarProcessorRequestSize(processor *ZBarProcessor, width, height uint32) int {
	return int(C.zbar_processor_request_size((*C.zbar_processor_t)(unsafe.Pointer(processor)), C.uint(width), C.uint(height)))
}

/** request a preferred video driver interface version for
 * debug/testing.
 * @note must be called before zbar_processor_init()
 * @since 0.6
 */
func ZBarProcessorRequestInterface(processor *ZBarProcessor, version int) int {
	return int(C.zbar_processor_request_interface((*C.zbar_processor_t)(unsafe.Pointer(processor)), C.int(version)))
}

/** request a preferred video I/O mode for debug/testing.  You will
 * get errors if the driver does not support the specified mode.
 * @verbatim
    0 = auto-detect
    1 = force I/O using read()
    2 = force memory mapped I/O using mmap()
    3 = force USERPTR I/O (v4l2 only)
@endverbatim
 * @note must be called before zbar_processor_init()
 * @since 0.7
 */
func ZBarProcessorRequestIomode(video *ZBarProcessor, iomode int) int {
	return int(C.zbar_processor_request_iomode((*C.zbar_processor_t)(unsafe.Pointer(video)), C.int(iomode)))
}

/** force specific input and output formats for debug/testing.
 * @note must be called before zbar_processor_init()
 */
func ZBarProcessorForceFormat(processor *ZBarProcessor, inputFormat, outputFormat uint64) int {
	return int(C.zbar_processor_force_format((*C.zbar_processor_t)(unsafe.Pointer(processor)), C.ulong(inputFormat), C.ulong(outputFormat)))
}

/** setup result handler callback.
 * the specified function will be called by the processor whenever
 * new results are available from the video stream or a static image.
 * pass a NULL value to disable callbacks.
 * @param processor the object on which to set the handler.
 * @param handler the function to call when new results are available.
 * @param userdata is set as with zbar_processor_set_userdata().
 * @returns the previously registered handler
 */
func ZBarProcessorSetDataHandler(processor *ZBarProcessor, handler ZBarImageDataHandler, userData unsafe.Pointer) ZBarImageDataHandler {
	// TODO...
	var handle = C.zbar_processor_set_data_handler((*C.zbar_processor_t)(unsafe.Pointer(processor)), (*C.zbar_image_data_handler_t)(unsafe.Pointer(syscall.NewCallback(handler))), userData)
	return func(image *ZBarImage, userData unsafe.Pointer){
		handle()
	}
}

/** associate user specified data value with the processor.
 * @since 0.6
 */
func ZBarProcessorSetUserData(processor *ZBarProcessor, userData unsafe.Pointer) {
	C.zbar_processor_set_userdata((*C.zbar_processor_t)(unsafe.Pointer(processor)), userData)
}

/** return user specified data value associated with the processor.
 * @since 0.6
 */
func ZBarProcessorGetUserData(processor *ZBarProcessor) unsafe.Pointer {
	return C.zbar_processor_get_userdata((*C.zbar_processor_t)(unsafe.Pointer(processor)))
}

/** set config for indicated symbology (0 for all) to specified value.
 * @returns 0 for success, non-0 for failure (config does not apply to
 * specified symbology, or value out of range)
 * @see zbar_decoder_set_config()
 * @since 0.4
 */
func ZBarProcessorSetConfig(processor *ZBarProcessor, symbology ZBarSymbolType, config ZBarConfig, value int) int {
	return int(C.zbar_processor_set_config((*C.zbar_processor_t)(unsafe.Pointer(processor)), C.zbar_symbol_type_t(symbology), C.zbar_config_t(config), C.int(value)))
}

/** parse configuration string using zbar_parse_config()
 * and apply to processor using zbar_processor_set_config().
 * @returns 0 for success, non-0 for failure
 * @see zbar_parse_config()
 * @see zbar_processor_set_config()
 * @since 0.4
 */
func ZBarProcessorParseConfig(processor *ZBarProcessor, configString string) int {
	var sym ZBarSymbolType
	var cfg ZBarConfig
	var val int

	// TODO
	if ret := ZBarParseConfig(configString, &sym, &cfg, &val); ret != 0 {
		return ret
	}

	return ZBarProcessorSetConfig(processor, sym, cfg, val)
}

/** retrieve the current state of the ouput window.
 * @returns 1 if the output window is currently displayed, 0 if not.
 * @returns -1 if an error occurs
 */
func ZBarProcessorIsVisible(processor *ZBarProcessor) int {
	return int(C.zbar_processor_is_visible((*C.zbar_processor_t)(unsafe.Pointer(processor))))
}

/** show or hide the display window owned by the library.
 * the size will be adjusted to the input size
 */
func ZBarProcessorSetVisible(processor *ZBarProcessor, visible int) int {
	return int(C.zbar_processor_set_visible((*C.zbar_processor_t)(unsafe.Pointer(processor)), C.int(visible)))
}

/** control the processor in free running video mode.
 * only works if video input is initialized. if threading is in use,
 * scanning will occur in the background, otherwise this is only
 * useful wrapping calls to zbar_processor_user_wait(). if the
 * library output window is visible, video display will be enabled.
 */
func ZBarProcessorSetActive(processor *ZBarProcessor, active int) int {
	return int(C.zbar_processor_set_active((*C.zbar_processor_t)(unsafe.Pointer(processor)), C.int(active)))
}

/** retrieve decode results for last scanned image/frame.
 * @returns the symbol set result container or NULL if no results are
 * available
 * @note the returned symbol set has its reference count incremented;
 * ensure that the count is decremented after use
 * @since 0.10
 */
func ZBarProcessorGetResults(processor *ZBarProcessor) *ZBarSymbolSet {
	return (*ZBarSymbolSet)(C.zbar_processor_get_results((*C.zbar_processor_t)(unsafe.Pointer(processor))))
}

/** wait for input to the display window from the user
 * (via mouse or keyboard).
 * @returns >0 when input is received, 0 if timeout ms expired
 * with no input or -1 in case of an error
 */
func ZBarProcessorUserWait(processor *ZBarProcessor, timeout int) int {
	return int(C.zbar_processor_user_wait((*C.zbar_processor_t)(unsafe.Pointer(processor)), C.int(timeout)))
}

/** process from the video stream until a result is available,
 * or the timeout (in milliseconds) expires.
 * specify a timeout of -1 to scan indefinitely
 * (zbar_processor_set_active() may still be used to abort the scan
 * from another thread).
 * if the library window is visible, video display will be enabled.
 * @note that multiple results may still be returned (despite the
 * name).
 * @returns >0 if symbols were successfully decoded,
 * 0 if no symbols were found (ie, the timeout expired)
 * or -1 if an error occurs
 */
func ZBarProcessOne(processor *ZBarProcessor, timeout int) int {
	return int(C.zbar_process_one((*C.zbar_processor_t)(unsafe.Pointer(processor)), C.int(timeout)))
}

/** process the provided image for barcodes.
 * if the library window is visible, the image will be displayed.
 * @returns >0 if symbols were successfully decoded,
 * 0 if no symbols were found or -1 if an error occurs
 */
func ZBarProcessImage(processor *ZBarProcessor, image *ZBarImage) int {
	return int(C.zbar_process_image((*C.zbar_processor_t)(unsafe.Pointer(processor)), (*C.zbar_image_t)(unsafe.Pointer(image))))
}

/** display detail for last processor error to stderr.
 * @returns a non-zero value suitable for passing to exit()
 */
func ZBarProcessorErrorSpew(processor *ZBarProcessor, verbosity int) int {
	return ZBarErrorSpew(processor, verbosity)
}

/** retrieve the detail string for the last processor error. */
func ZBarProcessorErrorString(processor *ZBarProcessor, verbosity int) string {
	return ZBarErrorString(processor, verbosity)
}

/** retrieve the type code for the last processor error. */
func ZBarProcessorGetErrorCode(processor *ZBarProcessor) ZBarError {
	return ZBarGetErrorCode(processor)
}

/*@}*/

/*------------------------------------------------------------*/
/** @name Video interface
 * @anchor c-video
 * mid-level video source abstraction.
 * captures images from a video device
 */
/*@{*/

type ZBarVideoStruct struct {}
/** opaque video object. */
type ZBarVideo ZBarVideoStruct

/** constructor. */
func ZBarVideoCreate() *ZBarVideo {
	return (*ZBarVideo)(C.zbar_video_create())
}

/** destructor. */
func ZBarVideoDestroy(video *ZBarVideo) {
	C.zbar_video_destroy((*C.zbar_video_t)(unsafe.Pointer(video)))
}

/** open and probe a video device.
 * the device specified by platform specific unique name
 * (v4l device node path in *nix eg "/dev/video",
 *  DirectShow DevicePath property in windows).
 * @returns 0 if successful or -1 if an error occurs
 */
func ZBarVideoOpen(video *ZBarVideo, device string) int {
	return int(C.zbar_video_open((*C.zbar_video_t)(unsafe.Pointer(video)), C.CString(device)))
}

/** retrieve file descriptor associated with open *nix video device
 * useful for using select()/poll() to tell when new images are
 * available (NB v4l2 only!!).
 * @returns the file descriptor or -1 if the video device is not open
 * or the driver only supports v4l1
 */
func ZBarVideoGetFd(video *ZBarVideo) int {
	return int(C.zbar_video_get_fd((*C.zbar_video_t)(unsafe.Pointer(video))))
}

/** request a preferred size for the video image from the device.
 * the request may be adjusted or completely ignored by the driver.
 * @returns 0 if successful or -1 if the video device is already
 * initialized
 * @since 0.6
 */
func ZBarVideoRequestSize(video *ZBarVideo, width, height uint32) int {
	return int(C.zbar_video_request_size((*C.zbar_video_t)(unsafe.Pointer(video)), C.uint(width), C.uint(height)))
}

/** request a preferred driver interface version for debug/testing.
 * @note must be called before zbar_video_open()
 * @since 0.6
 */
func ZBarVideoRequestInterface(video *ZBarVideo, version int) int {
	return int(C.zbar_video_request_interface((*C.zbar_video_t)(unsafe.Pointer(video)), C.int(version)))
}

/** request a preferred I/O mode for debug/testing.  You will get
 * errors if the driver does not support the specified mode.
 * @verbatim
    0 = auto-detect
    1 = force I/O using read()
    2 = force memory mapped I/O using mmap()
    3 = force USERPTR I/O (v4l2 only)
@endverbatim
 * @note must be called before zbar_video_open()
 * @since 0.7
 */
func ZBarVideoRequestIomode(video *ZBarVideo, iomode int) int {
	return int(C.zbar_video_request_iomode((*C.zbar_video_t)(unsafe.Pointer(video)), C.int(iomode)))
}

/** retrieve current output image width.
 * @returns the width or 0 if the video device is not open
 */
func ZBarVideoGetWidth(video *ZBarVideo) int {
	return int(C.zbar_video_get_width((*C.zbar_video_t)(unsafe.Pointer(video))))
}

/** retrieve current output image height.
 * @returns the height or 0 if the video device is not open
 */
func ZBarVideoGetHeight(video *ZBarVideo) int {
	return int(C.zbar_video_get_height((*C.zbar_video_t)(unsafe.Pointer(video))))
}

/** initialize video using a specific format for debug.
 * use zbar_negotiate_format() to automatically select and initialize
 * the best available format
 */
func ZBarVideoInit(video *ZBarVideo, format uint64) int {
	return int(C.zbar_video_init((*C.zbar_video_t)(unsafe.Pointer(video)), C.ulong(format)))
}

/** start/stop video capture.
 * all buffered images are retired when capture is disabled.
 * @returns 0 if successful or -1 if an error occurs
 */
func ZBarVideoEnable(video *ZBarVideo, enable int) int {
	return int(C.zbar_video_enable((*C.zbar_video_t)(unsafe.Pointer(video)), C.int(enable)))
}

/** retrieve next captured image.  blocks until an image is available.
 * @returns NULL if video is not enabled or an error occurs
 */
func ZBarVideoNextImage(video *ZBarVideo) *ZBarImage {
	return (*ZBarImage)(C.zbar_video_next_image((*C.zbar_video_t)(unsafe.Pointer(video))))
}

/** display detail for last video error to stderr.
 * @returns a non-zero value suitable for passing to exit()
 */
func ZBarVideoErrorSpew(video *ZBarVideo, verbosity int) int {
	return ZBarErrorSpew(video, verbosity)
}

/** retrieve the detail string for the last video error. */
func ZBarVideoErrorString(video *ZBarVideo, verbosity int) string {
	return ZBarErrorString(video, verbosity)
}

/** retrieve the type code for the last video error. */
func ZBarVideoGetErrorCode(video *ZBarVideo) ZBarError {
	return ZBarGetErrorCode(video)
}

/*@}*/

/*------------------------------------------------------------*/
/** @name Window interface
 * @anchor c-window
 * mid-level output window abstraction.
 * displays images to user-specified platform specific output window
 */
/*@{*/

type ZBarWindowStruct struct {}
/** opaque window object. */
type ZBarWindow ZBarWindowStruct

/** constructor. */
func ZBarWindowCreate() *ZBarWindow {
	return (*ZBarWindow)(C.zbar_window_create())
}

/** destructor. */
func ZBarWindowDestroy(window *ZBarWindow) {
	C.zbar_window_destroy((*C.zbar_window_t)(unsafe.Pointer(window)))
}

/** associate reader with an existing platform window.
 * This can be any "Drawable" for X Windows or a "HWND" for windows.
 * input images will be scaled into the output window.
 * pass NULL to detach from the resource, further input will be
 * ignored
 */
func ZBarWindowAttach(window *ZBarWindow, x11DisplayW32Hwnd unsafe.Pointer, x11Drawable uint64) int {
	return int(C.zbar_window_attach((*C.zbar_window_t)(unsafe.Pointer(window)), x11DisplayW32Hwnd, C.ulong(x11Drawable)))
}

/** control content level of the reader overlay.
 * the overlay displays graphical data for informational or debug
 * purposes.  higher values increase the level of annotation (possibly
 * decreasing performance). @verbatim
    0 = disable overlay
    1 = outline decoded symbols (default)
    2 = also track and display input frame rate
@endverbatim
 */
func ZBarWindowSetOverlay(window *ZBarWindow, level int) {
	C.zbar_window_set_overlay((*C.zbar_window_t)(unsafe.Pointer(window)), C.int(level))
}

/** retrieve current content level of reader overlay.
 * @see zbar_window_set_overlay()
 * @since 0.10
 */
func ZBarWindowGetOverlay(window *ZBarWindow) int {
	return int(C.zbar_window_get_overlay((*C.zbar_window_t)(unsafe.Pointer(window))))
}

/** draw a new image into the output window. */
func ZBarWindowDraw(window *ZBarWindow, image *ZBarImage) int {
	return int(C.zbar_window_draw((*C.zbar_window_t)(unsafe.Pointer(window)), (*C.zbar_image_t)(unsafe.Pointer(image))))
}

/** redraw the last image (exposure handler). */
func ZBarWindowRedraw(window *ZBarWindow) int {
	return int(C.zbar_window_redraw((*C.zbar_window_t)(unsafe.Pointer(window))))
}

/** resize the image window (reconfigure handler).
 * this does @em not update the contents of the window
 * @since 0.3, changed in 0.4 to not redraw window
 */
func ZBarWindowResize(window *ZBarWindow, width, height uint32) int {
	return int(C.zbar_window_redraw((*C.zbar_window_t)(unsafe.Pointer(window)), C.uint(width), C.uint(height)))
}

/** display detail for last window error to stderr.
 * @returns a non-zero value suitable for passing to exit()
 */
func ZBarWindowErrorSpew(window *ZBarWindow, verbosity int) int {
	return ZBarErrorSpew(window, verbosity)
}

/** retrieve the detail string for the last window error. */
func ZBarWindowErrorString(window *ZBarWindow, verbosity int) string {
	return ZBarErrorString(window, verbosity)
}

/** retrieve the type code for the last window error. */
func ZBarWindowGetErrorCode(window *ZBarWindow) ZBarError {
	return ZBarGetErrorCode(window)
}


/** select a compatible format between video input and output window.
 * the selection algorithm attempts to use a format shared by
 * video input and window output which is also most useful for
 * barcode scanning.  if a format conversion is necessary, it will
 * heuristically attempt to minimize the cost of the conversion
 */
func ZBarNegotiateFormat(video *ZBarVideo, window *ZBarWindow) int {
	return int(C.zbar_negotiate_format((*C.zbar_video_t)(unsafe.Pointer(video)), (*C.zbar_window_t)(unsafe.Pointer(window))))
}

/*@}*/

/*------------------------------------------------------------*/
/** @name Image Scanner interface
 * @anchor c-imagescanner
 * mid-level image scanner interface.
 * reads barcodes from 2-D images
 */
/*@{*/

type ZBarImageScannerStruct struct {}
/** opaque image scanner object. */
type ZBarImageScanner ZBarImageScannerStruct

/** constructor. */
func ZBarImageScannerCreate() *ZBarImageScanner {
	return (*ZBarImageScanner)(C.zbar_image_scanner_create())
}

/** destructor. */
func ZBarImageScannerDestroy(scanner *ZBarImageScanner) {
	C.zbar_image_scanner_destroy((*C.zbar_image_scanner_t)(unsafe.Pointer(scanner)))
}

/** setup result handler callback.
 * the specified function will be called by the scanner whenever
 * new results are available from a decoded image.
 * pass a NULL value to disable callbacks.
 * @returns the previously registered handler
 */
func ZBarImageScannerSetDataHandler(scanner *ZBarImageScanner, handler ZBarImageDataHandler, userData unsafe.Pointer) ZBarImageDataHandler {
	var fn = C.zbar_image_scanner_set_data_handler((*C.zbar_image_scanner_t)(unsafe.Pointer(scanner)), (*C.zbar_image_data_handler_t)(unsafe.Pointer(syscall.NewCallback(handler))), userData)
	// TODO ...
	return func(image *ZBarImage, userData unsafe.Pointer) {
		fn(image, userData)
	}
}


/** set config for indicated symbology (0 for all) to specified value.
 * @returns 0 for success, non-0 for failure (config does not apply to
 * specified symbology, or value out of range)
 * @see zbar_decoder_set_config()
 * @since 0.4
 */
func ZBarImageScannerSetConfig(scanner *ZBarImageScanner, symbology ZBarSymbolType, config ZBarConfig, value int) int {
	return int(C.zbar_image_scanner_set_config((*C.zbar_image_scanner_t)(unsafe.Pointer(scanner)), C.zbar_symbol_type_t(symbology), C.zbar_config_t(config), C.int(value)))
}

/** parse configuration string using zbar_parse_config()
 * and apply to image scanner using zbar_image_scanner_set_config().
 * @returns 0 for success, non-0 for failure
 * @see zbar_parse_config()
 * @see zbar_image_scanner_set_config()
 * @since 0.4
 */
func ZBarImageScannerParseConfig(scanner *ZBarImageScanner, configString string) int {
	var sym ZBarSymbolType
	var cfg ZBarConfig
	var val int

	if ret := ZBarParseConfig(configString, &sym, &cfg, &val); ret != 0 {
		return ret
	}

	return ZBarImageScannerSetConfig(scanner, sym, cfg, val)
}

/** enable or disable the inter-image result cache (default disabled).
 * mostly useful for scanning video frames, the cache filters
 * duplicate results from consecutive images, while adding some
 * consistency checking and hysteresis to the results.
 * this interface also clears the cache
 */
func ZBarImageScannerEnableCache(scanner *ZBarImageScanner, enable int) {
	C.zbar_image_scanner_enable_cache((*C.zbar_image_scanner_t)(unsafe.Pointer(scanner)), C.int(enable))
}

/** remove any previously decoded results from the image scanner and the
 * specified image.  somewhat more efficient version of
 * zbar_image_set_symbols(image, NULL) which may retain memory for
 * subsequent decodes
 * @since 0.10
 */
func ZBarImageScannerRecycleImage(scanner *ZBarImageScanner, image *ZBarImage) {
	C.zbar_image_scanner_recycle_image((*C.zbar_image_scanner_t)(unsafe.Pointer(scanner)), (*C.zbar_image_t)(unsafe.Pointer(image)))
}

/** retrieve decode results for last scanned image.
 * @returns the symbol set result container or NULL if no results are
 * available
 * @note the symbol set does not have its reference count adjusted;
 * ensure that the count is incremented if the results may be kept
 * after the next image is scanned
 * @since 0.10
 */
func ZBarImageScannerGetResults(scanner *ZBarImageScanner) *ZBarSymbolSet {
	return (*ZBarSymbolSet)(C.zbar_image_scanner_get_results((*C.zbar_image_scanner_t)(unsafe.Pointer(scanner))))
}

/** scan for symbols in provided image.  The image format must be
 * "Y800" or "GRAY".
 * @returns >0 if symbols were successfully decoded from the image,
 * 0 if no symbols were found or -1 if an error occurs
 * @see zbar_image_convert()
 * @since 0.9 - changed to only accept grayscale images
 */
func ZBarScanImage(scanner *ZBarImageScanner, image *ZBarImage) int {
	return int(C.zbar_scan_image((*C.zbar_image_scanner_t)(unsafe.Pointer(scanner)), (*C.zbar_image_t)(unsafe.Pointer(image))))
}

/*@}*/

/*------------------------------------------------------------*/
/** @name Decoder interface
 * @anchor c-decoder
 * low-level bar width stream decoder interface.
 * identifies symbols and extracts encoded data
 */
/*@{*/

type ZBarDecoderStruct struct {}
/** opaque decoder object. */
type ZBarDecoder ZBarDecoderStruct

/** decoder data handler callback function.
 * called by decoder when new data has just been decoded
 */
type ZBarDecoderHandler func(decoder *ZBarDecoder)

/** constructor. */
func ZBarDecoderCreate() *ZBarDecoder {
	return (*ZBarDecoder)(C.zbar_decoder_create())
}

/** destructor. */
func ZBarDecoderDestroy(decoder *ZBarDecoder) {
	C.zbar_decoder_destroy((*C.zbar_decoder_t)(unsafe.Pointer(decoder)))
}

/** set config for indicated symbology (0 for all) to specified value.
 * @returns 0 for success, non-0 for failure (config does not apply to
 * specified symbology, or value out of range)
 * @since 0.4
 */
func ZBarDecoderSetConfig(decoder *ZBarDecoder, symbology ZBarSymbolType, config ZBarConfig, value int) int {
	return int(C.zbar_decoder_set_config((*C.zbar_decoder_t)(unsafe.Pointer(decoder)), C.zbar_symbol_type_t(symbology), C.zbar_config_t(config), C.int(value)))
}

/** parse configuration string using zbar_parse_config()
 * and apply to decoder using zbar_decoder_set_config().
 * @returns 0 for success, non-0 for failure
 * @see zbar_parse_config()
 * @see zbar_decoder_set_config()
 * @since 0.4
 */
func zBarDecoderParseConfig(decoder *ZBarDecoder, configString string) int {
	var sym ZBarSymbolType
	var cfg ZBarConfig
	var val int

	if ret := ZBarParseConfig(configString, &sym, &cfg, &val); ret != 0 {
		return ret
	}

	return ZBarDecoderSetConfig(decoder, sym, cfg, val)
}

/** clear all decoder state.
 * any partial symbols are flushed
 */
func ZBarDecoderReset(decoder *ZBarDecoder) {
	C.zbar_decoder_reset((*C.zbar_decoder_t)(unsafe.Pointer(decoder)))
}

/** mark start of a new scan pass.
 * clears any intra-symbol state and resets color to ::ZBAR_SPACE.
 * any partially decoded symbol state is retained
 */
func ZBarDecoderNewScan(decoder *ZBarDecoder) {
	C.zbar_decoder_new_scan((*C.zbar_decoder_t)(unsafe.Pointer(decoder)))
}

/** process next bar/space width from input stream.
 * the width is in arbitrary relative units.  first value of a scan
 * is ::ZBAR_SPACE width, alternating from there.
 * @returns appropriate symbol type if width completes
 * decode of a symbol (data is available for retrieval)
 * @returns ::ZBAR_PARTIAL as a hint if part of a symbol was decoded
 * @returns ::ZBAR_NONE (0) if no new symbol data is available
 */
func ZBarDecodeWidth(decoder *ZBarDecoder, width uint32) ZBarSymbolType {
	return ZBarSymbolType(C.zbar_decode_width((*C.zbar_decoder_t)(unsafe.Pointer(decoder)), C.uint(width)))
}

/** retrieve color of @em next element passed to
 * zbar_decode_width(). */
func ZBarDecoderGetColor(decoder *ZBarDecoder) ZBarColor {
	return ZBarColor(C.zbar_decoder_get_color((*C.zbar_decoder_t)(unsafe.Pointer(decoder))))
}

/** retrieve last decoded data.
 * @returns the data string or NULL if no new data available.
 * the returned data buffer is owned by library, contents are only
 * valid between non-0 return from zbar_decode_width and next library
 * call
 */
func ZBarDecoderGetData(decoder *ZBarDecoder) string {
	return C.GoString(C.zbar_decoder_get_data((*C.zbar_decoder_t)(unsafe.Pointer(decoder))))
}

/** retrieve length of binary data.
 * @returns the length of the decoded data or 0 if no new data
 * available.
 */
func ZBarDecoderGetDataLength(decoder *ZBarDecoder) uint32 {
	return uint32(C.zbar_decoder_get_data_length((*C.zbar_decoder_t)(unsafe.Pointer(decoder))))
}

/** retrieve last decoded symbol type.
 * @returns the type or ::ZBAR_NONE if no new data available
 */
func ZBarDecoderGetType(decoder *ZBarDecoder) ZBarSymbolType {
	return ZBarSymbolType(C.zbar_decoder_get_type((*C.zbar_decoder_t)(unsafe.Pointer(decoder))))
}

/** setup data handler callback.
 * the registered function will be called by the decoder
 * just before zbar_decode_width() returns a non-zero value.
 * pass a NULL value to disable callbacks.
 * @returns the previously registered handler
 */
func ZBarDecoderSetHandler(decoder *ZBarDecoder, handler ZBarDecoderHandler) ZBarDecoderHandler {
	var fn = C.zbar_decoder_set_handler((*C.zbar_decoder_t)(unsafe.Pointer(decoder)), (*C.zbar_decoder_handler_t)(unsafe.Pointer(syscall.NewCallback(handler))))
	// TODO...
	return func(decoder *ZBarDecoder) {
		fn((*C.zbar_decoder_t)(unsafe.Pointer(decoder)))
	}
}

/** associate user specified data value with the decoder. */
func ZBarDecoderSetUserData(decoder *ZBarDecoder, userData unsafe.Pointer) {
	C.zbar_decoder_set_userdata((*C.zbar_decoder_t)(unsafe.Pointer(decoder)), userData)
}

/** return user specified data value associated with the decoder. */
func ZBarDecoderGetUserData(decoder *ZBarDecoder) unsafe.Pointer {
	return C.zbar_decoder_get_userdata((*C.zbar_decoder_t)(unsafe.Pointer(decoder)))
}

/*@}*/

/*------------------------------------------------------------*/
/** @name Scanner interface
 * @anchor c-scanner
 * low-level linear intensity sample stream scanner interface.
 * identifies "bar" edges and measures width between them.
 * optionally passes to bar width decoder
 */
/*@{*/

type ZBarScannerStruct struct {}
/** opaque scanner object. */
type ZBarScanner ZBarScannerStruct

/** constructor.
 * if decoder is non-NULL it will be attached to scanner
 * and called automatically at each new edge
 * current color is initialized to ::ZBAR_SPACE
 * (so an initial BAR->SPACE transition may be discarded)
 */
func ZBarScannerCreate(decoder *ZBarDecoder) *ZBarScanner {
	return (*ZBarScanner)(C.zbar_scanner_create())
}

/** destructor. */
func ZBarScannerDestroy(scanner *ZBarScanner) {
	C.zbar_scanner_destroy((*C.zbar_scanner_t)(unsafe.Pointer(scanner)))
}

/** clear all scanner state.
 * also resets an associated decoder
 */
func ZBarScannerReset(scanner *ZBarScanner) ZBarSymbolType {
	return ZBarSymbolType(C.zbar_scanner_reset((*C.zbar_scanner_t)(unsafe.Pointer(scanner))))
}

/** mark start of a new scan pass. resets color to ::ZBAR_SPACE.
 * also updates an associated decoder.
 * @returns any decode results flushed from the pipeline
 * @note when not using callback handlers, the return value should
 * be checked the same as zbar_scan_y()
 * @note call zbar_scanner_flush() at least twice before calling this
 * method to ensure no decode results are lost
 */
func ZBarScannerNewScan(scanner *ZBarScanner) ZBarSymbolType {
	return ZBarSymbolType(C.zbar_scanner_new_scan((*C.zbar_scanner_t)(unsafe.Pointer(scanner))))
}

/** flush scanner processing pipeline.
 * forces current scanner position to be a scan boundary.
 * call multiple times (max 3) to completely flush decoder.
 * @returns any decode/scan results flushed from the pipeline
 * @note when not using callback handlers, the return value should
 * be checked the same as zbar_scan_y()
 * @since 0.9
 */
func ZBarScannerFlush(scanner *ZBarScanner) ZBarSymbolType {
	return ZBarSymbolType(C.zbar_scanner_flush((*C.zbar_scanner_t)(unsafe.Pointer(scanner))))
}

/** process next sample intensity value.
 * intensity (y) is in arbitrary relative units.
 * @returns result of zbar_decode_width() if a decoder is attached,
 * otherwise @returns (::ZBAR_PARTIAL) when new edge is detected
 * or 0 (::ZBAR_NONE) if no new edge is detected
 */
func ZBarScanY(scanner *ZBarScanner, y int) ZBarSymbolType {
	return ZBarSymbolType(C.zbar_scan_y((*C.zbar_scanner_t)(unsafe.Pointer(scanner)), C.int(y)))
}

/** process next sample from RGB (or BGR) triple. */
func ZBarScanRgb24(scanner *ZBarScanner, rgb *byte) ZBarSymbolType {
	return ZBarSymbolType(C.zbar_scan_rgb24((*C.zbar_scanner_t)(unsafe.Pointer(scanner)), (*C.uchar)(unsafe.Pointer(rgb))))
}

/** retrieve last scanned width. */
func ZBarScannerGetWidth(scanner *ZBarScanner) uint32 {
	return uint32(C.zbar_scanner_get_width((*C.zbar_scanner_t)(unsafe.Pointer(scanner))))
}

/** retrieve sample position of last edge.
 * @since 0.10
 */
func ZBarScannerGetEdge(scn *ZBarScanner, offset uint32, prec int) uint32 {
	return uint32(C.zbar_scanner_get_edge((*C.zbar_scanner_t)(unsafe.Pointer(scn)), C.uint(offset), C.int(prec)))
}

/** retrieve last scanned color. */
func ZBarScannerGetColor(scanner *ZBarScanner) ZBarColor {
	return ZBarColor(C.zbar_scanner_get_color((*C.zbar_scanner_t)(unsafe.Pointer(scanner))))
}

/*@}*/




















func Type(t interface{}) {
	fmt.Println(t)
	fmt.Println(reflect.TypeOf(t).String())
}

func test() {
	Type(ZBAR_CFG_MAX_LEN)
	Type(ZBAR_CFG_POSITION)

	fmt.Println("===================")
	fmt.Println(ZBAR_CFG_ENABLE )
	fmt.Println(ZBAR_CFG_ADD_CHECK )
	fmt.Println(ZBAR_CFG_EMIT_CHECK )
	fmt.Println(ZBAR_CFG_ASCII )
	fmt.Println(ZBAR_CFG_NUM )
	fmt.Println(ZBAR_CFG_MIN_LEN )
	fmt.Println(ZBAR_CFG_MAX_LEN )
	fmt.Println(ZBAR_CFG_POSITION )
	fmt.Println(ZBAR_CFG_X_DENSITY )
	fmt.Println(ZBAR_CFG_Y_DENSITY)

	fmt.Println("========================")
	fmt.Println(ZBarGetSymbolName(ZBAR_NONE))
}

func main() {
	var major,minor uint32
	fmt.Println("ZBarVersion:", ZBarVersion(&major, &minor))
	fmt.Println("major:", major)
	fmt.Println("minor:", minor)
	ZBarSetVerbosity(100)
	ZBarIncreaseVerbosity()
	fmt.Println("============================")
	fmt.Println(ZBarGetSymbolName(ZBAR_NONE))
	fmt.Println(ZBarGetSymbolName(ZBAR_PARTIAL))
	fmt.Println(ZBarGetSymbolName(ZBAR_EAN8))
	fmt.Println(ZBarGetSymbolName(ZBAR_UPCE))
	fmt.Println(ZBarGetSymbolName(ZBAR_ISBN10))
	fmt.Println(ZBarGetSymbolName(ZBAR_UPCA))
	fmt.Println(ZBarGetSymbolName(ZBAR_EAN13))
	fmt.Println(ZBarGetSymbolName(ZBAR_ISBN13))
	fmt.Println(ZBarGetSymbolName(ZBAR_I25))
	fmt.Println(ZBarGetSymbolName(ZBAR_CODE39))
	fmt.Println(ZBarGetSymbolName(ZBAR_PDF417))
	fmt.Println(ZBarGetSymbolName(ZBAR_QRCODE))
	fmt.Println(ZBarGetSymbolName(ZBAR_CODE128))
	fmt.Println(ZBarGetSymbolName(ZBAR_SYMBOL))
	fmt.Println(ZBarGetSymbolName(ZBAR_ADDON2))
	fmt.Println(ZBarGetSymbolName(ZBAR_ADDON5))
	fmt.Println(ZBarGetSymbolName(ZBAR_ADDON))
	fmt.Println("============================")
	fmt.Println(ZBarGetAddonName(ZBAR_NONE))
	fmt.Println(ZBarGetAddonName(ZBAR_PARTIAL))
	fmt.Println(ZBarGetAddonName(ZBAR_EAN8))
	fmt.Println(ZBarGetAddonName(ZBAR_UPCE))
	fmt.Println(ZBarGetAddonName(ZBAR_ISBN10))
	fmt.Println(ZBarGetAddonName(ZBAR_UPCA))
	fmt.Println(ZBarGetAddonName(ZBAR_EAN13))
	fmt.Println(ZBarGetAddonName(ZBAR_ISBN13))
	fmt.Println(ZBarGetAddonName(ZBAR_I25))
	fmt.Println(ZBarGetAddonName(ZBAR_CODE39))
	fmt.Println(ZBarGetAddonName(ZBAR_PDF417))
	fmt.Println(ZBarGetAddonName(ZBAR_QRCODE))
	fmt.Println(ZBarGetAddonName(ZBAR_CODE128))
	fmt.Println(ZBarGetAddonName(ZBAR_SYMBOL))
	fmt.Println(ZBarGetAddonName(ZBAR_ADDON2))
	fmt.Println(ZBarGetAddonName(ZBAR_ADDON5))
	fmt.Println(ZBarGetAddonName(ZBAR_ADDON))
	fmt.Println("===============================")
	var symbology ZBarSymbolType
	var config ZBarConfig
	var value int
	fmt.Println(ZBarParseConfig("", &symbology, &config, &value))
	fmt.Println("symbology:", symbology)
	fmt.Println("config:", config)
	fmt.Println("value:", value)
	//fmt.Println("ZBarErrorSpew", ZBarErrorSpew(nil, 1))
	//fmt.Println("ZBarErrorSpew", ZBarErrorString(nil, 2))
	//fmt.Println("ZBarGetErrorCode", ZBarGetErrorCode(nil))
	var symbol ZBarSymbol
	ZBarSymbolRef(&symbol, 2)
	fmt.Println(symbol)
	fmt.Println(ZBarSymbolGetType(&symbol))
}