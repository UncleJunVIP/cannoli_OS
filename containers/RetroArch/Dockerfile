FROM crosstool-ng-builder:latest AS toolchain-base

USER root

RUN apt-get update && apt-get install -y \
    pkg-config libasound2-dev \
    libsdl2-dev libxkbcommon-dev \
    && rm -rf /var/lib/apt/lists/*

USER builder
WORKDIR /home/builder

RUN git clone https://github.com/libretro/RetroArch.git

WORKDIR /home/builder/RetroArch

RUN  ./configure \
        --enable-sdl2 \
        --enable-alsa \
        --disable-pulse \
        --disable-opengl \
        --disable-opengles \
        --disable-kms \
        --disable-x11 \
        --disable-wayland \
        --disable-vulkan \
        --disable-xvideo \
        --disable-xinerama \
        --disable-xshm \
        --disable-xrandr \
        --disable-ffmpeg

RUN make clean
RUN make -j$(nproc)

WORKDIR /home/builder

RUN git clone https://github.com/libretro/gambatte-libretro.git

WORKDIR /home/builder/gambatte-libretro

RUN make -f Makefile.libretro -j$(nproc)

RUN git clone https://github.com/libretro/libretro-core-info.git

RUN mkdir -p /home/builder/RetroArch/info
RUN cp libretro-core-info/gambatte_libretro.info /home/builder/RetroArch/info/

ENTRYPOINT ["bash"]