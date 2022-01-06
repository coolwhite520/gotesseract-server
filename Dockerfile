##### docker build -f docker_file -t otiai10:0.1 .
# This is a working example of setting up tesseract/gosseract,
# and also works as an example runtime to use gosseract package.
# You can just hit `docker run -it --rm otiai10/gosseract`
# to try and check it out!

#####
FROM golang:latest

RUN apt-get update -qq
# You need librariy files and headers of tesseract and leptonica.
# When you miss these or LD_LIBRARY_PATH is not set to them,
# you would face an error: "tesseract/baseapi.h: No such file or directory"
RUN apt-get install -y -qq libtesseract-dev libleptonica-dev
# In case you face TESSDATA_PREFIX error, you minght need to set env vars
# to specify the directory where "tessdata" is located.
ENV TESSDATA_PREFIX=/usr/share/tesseract-ocr/4.00/tessdata/

# Load languages.
# These {lang}.traineddata would b located under ${TESSDATA_PREFIX}/tessdata.
# 安装中文
#RUN apt-get install -y -qq \
#  tesseract-ocr-ara \
#  tesseract-ocr-deu \
#  tesseract-ocr-dan \
#  tesseract-ocr-eng \
#  tesseract-ocr-fra \
#  tesseract-ocr-mal \
#  tesseract-ocr-mom \
#  tesseract-ocr-tha \
#  tesseract-ocr-rus \
#  tesseract-ocr-jpn \
#  tesseract-ocr-chi-sim \
#  tesseract-ocr-vie \
#  tesseract-ocr-mya \
#  tesseract-ocr-nep \
#  tesseract-ocr-ita \
#  tesseract-ocr-kor
 #检查「tesseract」支持的语言
 # tesseract --list-langs
# See https://github.com/tesseract-ocr/tessdata for the list of available languages.
# https://github.com/tesseract-ocr/tessdata/tree/4.00 或者下载 https://tesseract-ocr.github.io/tessdoc/Data-Files
# If you want to download these traineddata via `wget`, don't forget to locate
# downloaded traineddata under ${TESSDATA_PREFIX}/tessdata.

COPY ./tessdata/* ${TESSDATA_PREFIX}

WORKDIR /go/src/app
ADD ./main.go /go/src/app
ADD ./go.mod /go/src/app
RUN go env -w GO111MODULE=on
RUN go mod tidy
EXPOSE 9090
RUN go build
RUN ["rm", "-rf", "/go/src/app/main.go", "/go/src/app/go.mod"]
CMD /go/src/app/gosseract
#RUN /go/src
#RUN cd ${GOPATH}/src/github.com/otiai10/gosseract && go test

# Now, you've got complete environment to play with "gosseract"!
# For other OS, check https://github.com/otiai10/gosseract/tree/main/test/runtimes