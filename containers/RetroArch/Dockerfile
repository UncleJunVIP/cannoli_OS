FROM crosstool-ng-builder:latest AS toolchain-base

USER root

RUN apt-get update && apt-get install -y \
    pkg-config libasound2-dev \
    libsdl2-dev libxkbcommon-dev \
    cmake \
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

RUN mkdir -p /home/builder/out
RUN mkdir -p /home/builder/out/cores
RUN mkdir -p /home/builder/out/cores/info
RUN mkdir -p /home/builder/out/lib

RUN cp /home/builder/RetroArch/retroarch /home/builder/out

WORKDIR /home/builder

# Grab Core Info
RUN git clone https://github.com/libretro/libretro-core-info.git
RUN cp libretro-core-info/*.info /home/builder/out/cores/info

# Clone Cores
RUN git clone https://github.com/libretro/gambatte-libretro.git
RUN git clone https://github.com/libretro/mgba.git
RUN git clone https://github.com/libretro/nestopia.git
RUN git clone https://github.com/libretro/snes9x.git
RUN git clone https://github.com/libretro/Genesis-Plus-GX.git
RUN git clone https://github.com/libretro/swanstation.git

# Build Gambatte (Game Boy/Game Boy Color)
WORKDIR /home/builder/gambatte-libretro

RUN make -f Makefile.libretro -j$(nproc)
RUN cp /home/builder/gambatte-libretro/gambatte_libretro.so /home/builder/out/cores

WORKDIR /home/builder

# Build mGBA (Game Boy Advance)
WORKDIR /home/builder/mgba

RUN make -f Makefile.libretro -j$(nproc)
RUN cp /home/builder/mgba/mgba_libretro.so /home/builder/out/cores

# Build Nestopia (Nintendo Entertainment System)
WORKDIR /home/builder/nestopia/libretro

RUN make -j$(nproc)
RUN cp /home/builder/nestopia/libretro/nestopia_libretro.so /home/builder/out/cores

# Build Snes9x (Super Nintendo Entertainment System)
WORKDIR /home/builder/snes9x/libretro

RUN make -j$(nproc)
RUN cp /home/builder/snes9x/libretro/snes9x_libretro.so /home/builder/out/cores

# Build Genesis Plus GX (Sega Genesis)
WORKDIR /home/builder/Genesis-Plus-GX

RUN make -f Makefile.libretro -j$(nproc)
RUN cp /home/builder/Genesis-Plus-GX/genesis_plus_gx_libretro.so /home/builder/out/cores

# Build SwanStation (PlayStation)
WORKDIR /home/builder/swanstation

RUN cmake -DCMAKE_BUILD_TYPE=Release -DBUILD_LIBRETRO_CORE=ON .
RUN make -j$(nproc)
RUN cp /home/builder/swanstation/swanstation_libretro.so /home/builder/out/cores

RUN cp /usr/lib/*/libxkbcommon.so.0 /home/builder/out/lib/

WORKDIR /home/builder

# Build FCEUmm (Nintendo Entertainment System)
RUN git clone https://github.com/libretro/libretro-fceumm

WORKDIR /home/builder/libretro-fceumm

RUN ls

RUN make -f Makefile.libretro -j$(nproc)
RUN cp /home/builder/libretro-fceumm/fceumm_libretro.so /home/builder/out/cores

WORKDIR /home/builder

ENTRYPOINT ["bash"]
