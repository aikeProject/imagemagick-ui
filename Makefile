GOCMD = go
GOBUILD = $(GOCMD) build
GOMOD = $(GOCMD) mod
GOTEST = $(GOCMD) test
# 指定ImageMagick文件所在位置
# Go语言中调用c 通过 CGO_CFLAGS CGO_LDFLAGS 环境变量，指定其所在位置
CGO_CFLAGS_IMAGICK="-Xpreprocessor -fopenmp -DMAGICKCORE_HDRI_ENABLE=1 -DMAGICKCORE_QUANTUM_DEPTH=16 -I/Users/chengyu/Downloads/ImageMagick-7.0.10/include/ImageMagick-7"
CGO_LDFLAGS_IMAGICK="-g -O2 -L/Users/chengyu/Downloads/ImageMagick-7.0.10/lib -lMagickWand-7.Q16HDRI -lMagickCore-7.Q16HDRI"

init:
	$(GOMOD) init $(module)

install:
	$(GOMOD) tidy

serve:
	CGO_CFLAGS=$(CGO_CFLAGS_IMAGICK) CGO_LDFLAGS=$(CGO_LDFLAGS_IMAGICK) $(GOBUILD) -v -x -tags no_pkgconfig gopkg.in/gographics/imagick.v3/imagick
	wails serve

build:
	CGO_CFLAGS=$(CGO_CFLAGS_IMAGICK) CGO_LDFLAGS=$(CGO_LDFLAGS_IMAGICK) $(GOBUILD) -v -x -tags no_pkgconfig gopkg.in/gographics/imagick.v3/imagick
	wails build -p

build-debug:
	wails build -d

# 清空go语言编译时的缓存
# rm -rfv ~/Library/Caches/go-build/