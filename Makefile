AppName=imagemagick-ui
NameMacApp=$(AppName).app
ImageMagick=ImageMagick-7.0.10
Frameworks=build/$(NameMacApp)/Contents/Frameworks
AppMacOS=build/$(NameMacApp)/Contents/MacOS/$(AppName)
LibMagickWand=libMagickWand-7.Q16HDRI.8.dylib
LibMagickCore=libMagickCore-7.Q16HDRI.8.dylib
LibMagickWandFile=/tmp/$(ImageMagick)/lib/$(LibMagickWand)
LibMagickCoreFile=/tmp/$(ImageMagick)/lib/$(LibMagickCore)
GOCMD = go
GOBUILD = $(GOCMD) build
GOMOD = $(GOCMD) mod
GOTEST = $(GOCMD) test
TOOl_ID = install_name_tool -id
TOOl_CHANGE = install_name_tool -change
# 指定ImageMagick文件所在位置
# Go语言中调用c 通过 CGO_CFLAGS CGO_LDFLAGS 环境变量，指定其所在位置
CGO_CFLAGS_IMAGICK=-Xpreprocessor -fopenmp -DMAGICKCORE_HDRI_ENABLE=1 -DMAGICKCORE_QUANTUM_DEPTH=16 -I/tmp/$(ImageMagick)/include/ImageMagick-7
CGO_LDFLAGS_IMAGICK=-g -O2 -L/tmp/$(ImageMagick)/lib/ -lMagickWand-7.Q16HDRI.8 -lMagickCore-7.Q16HDRI.8

# 全局环境变量
export CGO_CFLAGS=$(CGO_CFLAGS_IMAGICK)
export CGO_LDFLAGS=$(CGO_LDFLAGS_IMAGICK)

init:
	rm -rfv /tmp/$(ImageMagick)
	cp -rpv source/$(ImageMagick) /tmp/$(ImageMagick)
	$(TOOl_ID) $(LibMagickWandFile) $(LibMagickWandFile)
	$(TOOl_CHANGE) /$(ImageMagick)/lib/$(LibMagickCore) $(LibMagickCoreFile) $(LibMagickWandFile)
	$(TOOl_ID) $(LibMagickCoreFile) $(LibMagickCoreFile)

install:
	$(GOMOD) tidy

# 局部环境变量
#serve: export CGO_CFLAGS=$(CGO_CFLAGS_IMAGICK)
#serve: export CGO_LDFLAGS=$(CGO_LDFLAGS_IMAGICK)
serve:
	$(GOBUILD) -v -x -tags no_pkgconfig gopkg.in/gographics/imagick.v3/imagick
	wails serve

build:
	rm -rfv build/$(NameMacApp)
	$(GOBUILD) -v -x -tags no_pkgconfig gopkg.in/gographics/imagick.v3/imagick
	wails build -p
	mkdir -p $(Frameworks)
	cp -rpv /tmp/$(ImageMagick)/lib/*.8.dylib $(Frameworks)
	# 修改"ImageMagick"三方库动态链接地址
	# 添加"@rpath"地址
	install_name_tool -add_rpath @loader_path/../Frameworks $(AppMacOS)
	$(TOOl_CHANGE) $(LibMagickWandFile) @rpath/$(LibMagickWand) $(AppMacOS)
	$(TOOl_CHANGE) $(LibMagickCoreFile) @rpath/$(LibMagickCore) $(AppMacOS)
	$(TOOl_ID) @rpath/$(LibMagickWand) $(Frameworks)/$(LibMagickWand)
	$(TOOl_CHANGE) $(LibMagickCoreFile) @rpath/$(LibMagickCore) $(Frameworks)/$(LibMagickWand)
	$(TOOl_ID) @rpath/$(LibMagickCore) $(Frameworks)/$(LibMagickCore)

build-debug:
	$(GOBUILD) -v -x -tags no_pkgconfig gopkg.in/gographics/imagick.v3/imagick
	wails build -d

.PHONY: build

# 清空go语言编译时的缓存
# rm -rfv ~/Library/Caches/go-build/