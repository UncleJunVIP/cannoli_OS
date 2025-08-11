FROM debian:12-slim

ARG DEBIAN_FRONTEND=noninteractive
ARG CORES="2048 fceumm"
ARG RETROARCH_REPO=https://github.com/libretro/RetroArch.git
ARG RETROARCH_DEPTH=1
ARG CORE_2048_REPO=https://github.com/libretro/libretro-2048.git
ARG CORE_FCEUMM_REPO=https://github.com/libretro/libretro-fceumm.git
ARG ASSETS_REPO=https://github.com/libretro/retroarch-assets.git
ARG AUTOCONF_REPO=https://github.com/libretro/retroarch-joypad-autoconfig.git

# Enable arm64 multiarch and install toolchain + arm64 dev/runtime libs needed for features
RUN dpkg --add-architecture arm64 && apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates curl git file xz-utils bzip2 findutils \
    build-essential autoconf automake libtool cmake make pkg-config python3 \
    gcc-aarch64-linux-gnu g++-aarch64-linux-gnu binutils-aarch64-linux-gnu \
    libc6-dev-arm64-cross linux-libc-dev:arm64 \
    # Core compression/image/font deps
    zlib1g:arm64 zlib1g-dev:arm64 \
    libpng16-16:arm64 libpng-dev:arm64 \
    libjpeg62-turbo:arm64 libjpeg62-turbo-dev:arm64 \
    libfreetype6:arm64 libfreetype6-dev:arm64 \
    # Audio
    libasound2:arm64 libasound2-dev:arm64 \
    # Input (udev/evdev)
    libudev1:arm64 libudev-dev:arm64 libevdev2:arm64 libevdev-dev:arm64 \
    # SDL2 (headers/libs for building RetroArch; runtime SDL will be our custom build)
    libsdl2-2.0-0:arm64 libsdl2-dev:arm64 \
    # X11 bits (left available but RetroArch built with --disable-x11)
    libx11-6:arm64 libx11-dev:arm64 libxext6:arm64 libxext-dev:arm64 \
    libxrandr2:arm64 libxrandr-dev:arm64 libxi6:arm64 libxi-dev:arm64 libxinerama1:arm64 libxinerama-dev:arm64 \
    # GL/EGL/GLES
    libgl1-mesa-dri:arm64 libgl1-mesa-dev:arm64 \
    libegl1:arm64 libegl1-mesa:arm64 libegl1-mesa-dev:arm64 \
    libgles2:arm64 libgles2-mesa-dev:arm64 \
    # KMS/DRM/GBM (for EGL paths)
    libdrm2:arm64 libdrm-dev:arm64 libgbm1:arm64 libgbm-dev:arm64 \
    # Vulkan loader/headers
    libvulkan1:arm64 libvulkan-dev:arm64 \
    # glslang headers (arch-independent) and tools (arm64)
    glslang-dev glslang-tools:arm64 \
    # Optional SSL for netplay
    libssl3:arm64 libssl-dev:arm64 \
 && rm -rf /var/lib/apt/lists/*

# Auto-detect and install the correct libglslang runtime for arm64 (name varies by release)
RUN set -eux; \
  apt-get update; \
  GLSLANG_LIB="$(apt-cache search -n '^libglslang[0-9]+$' | awk '{print $1}' | head -n1 || true)"; \
  if [ -n "$GLSLANG_LIB" ]; then \
    echo "Installing ${GLSLANG_LIB}:arm64"; \
    apt-get install -y --no-install-recommends "${GLSLANG_LIB}:arm64"; \
  else \
    echo "WARN: Could not find a libglslangXX package; glslang features may be reduced."; \
  fi; \
  rm -rf /var/lib/apt/lists/*

WORKDIR /src

# Cross settings
ENV TARGET=aarch64-linux-gnu
ENV APPDIR=/opt/appdir
ENV OUTDIR=/out
RUN mkdir -p "$APPDIR" "$OUTDIR"

# Make pkg-config resolve arm64 .pc files
ENV PKG_CONFIG_DIR=
ENV PKG_CONFIG_SYSROOT_DIR=/
ENV PKG_CONFIG_LIBDIR=/usr/lib/aarch64-linux-gnu/pkgconfig

# Common cross variables
ENV CC=$TARGET-gcc \
    CXX=$TARGET-g++ \
    AR=$TARGET-gcc-ar \
    RANLIB=$TARGET-gcc-ranlib \
    STRIP=$TARGET-strip \
    CFLAGS="-O2 -pipe -fPIC -march=armv8-a" \
    CXXFLAGS="-O2 -pipe -fPIC -march=armv8-a"

SHELL ["/bin/bash", "-o", "pipefail", "-c"]

# Build RetroArch with SDL2, UDEV/EVDEV, OpenGLES, EGL, Vulkan, ALSA, OSS, LibretroDB, Netplay, glslang
RUN set -eux; \
  rm -rf RetroArch && git clone --depth="${RETROARCH_DEPTH}" "${RETROARCH_REPO}" RetroArch; \
  cd RetroArch; \
  git clean -fdx || true; \
  PKG_CONFIG_LIBDIR=/usr/lib/aarch64-linux-gnu/pkgconfig ./configure --host="${TARGET}" \
    --enable-sdl2 \
    --enable-opengles \
    --enable-egl \
    --enable-vulkan \
    --enable-udev \
    --enable-alsa \
    --enable-oss \
    --disable-opengl \
    --disable-x11; \
  make -j"$(nproc)" \
    CC="${TARGET}-gcc" CXX="${TARGET}-g++" HOSTCC=gcc \
    HAVE_THREADS=1 \
    HAVE_SDL2=1 \
    HAVE_UDEV=1 \
    HAVE_OPENGL=0 \
    HAVE_OPENGLES=1 \
    HAVE_EGL=1 \
    HAVE_VULKAN=1 \
    HAVE_ALSA=1 \
    HAVE_OSS=1 \
    HAVE_ZLIB=1 \
    HAVE_7ZIP=1 \
    HAVE_LIBRETRODB=1 \
    HAVE_NETWORK_CMD=1 \
    HAVE_COMMAND=1 \
    HAVE_GLSLANG=1 \
    HAVE_GLSL=1 \
    HAVE_DYNAMIC=1 \
    HAVE_X11=0; \
  "${STRIP}" retroarch || true; \
  ./retroarch --features || true

# Build a couple of simple cores (add more as needed)
RUN set -eux; mkdir -p /src/cores-built; \
  for core in $CORES; do \
    if [[ "$core" == "2048" ]]; then \
      git clone --depth=1 "$CORE_2048_REPO" core-2048; \
      make -C core-2048 -f Makefile.libretro -j"$(nproc)" \
        platform=unix \
        CC=$TARGET-gcc CXX=$TARGET-g++ \
        AR=$TARGET-gcc-ar RANLIB=$TARGET-gcc-ranlib STRIP=$TARGET-strip; \
      $STRIP core-2048/2048_libretro.so || true; \
      cp core-2048/2048_libretro.so /src/cores-built/; \
      rm -rf core-2048; \
    elif [[ "$core" == "fceumm" ]]; then \
      git clone --depth=1 "$CORE_FCEUMM_REPO" core-fceumm; \
      make -C core-fceumm -f Makefile.libretro -j"$(nproc)" \
        platform=unix \
        CC=$TARGET-gcc CXX=$TARGET-g++ \
        AR=$TARGET-gcc-ar RANLIB=$TARGET-gcc-ranlib STRIP=$TARGET-strip; \
      $STRIP core-fceumm/fceumm_libretro.so || true; \
      cp core-fceumm/fceumm_libretro.so /src/cores-built/; \
      rm -rf core-fceumm; \
    fi; \
  done

# Assets and autoconfigs
RUN set -eux; \
  git clone --depth=1 "$ASSETS_REPO" retroarch-assets; \
  git clone --depth=1 "$AUTOCONF_REPO" retroarch-joypad-autoconfig

# Assemble the portable appdir layout and copy RetroArch + assets + cores
RUN set -eux; \
  mkdir -p "$APPDIR"/{lib,usr/lib,cores,assets,autoconfig,system,shaders,records,playlists,downloads}; \
  cp /src/RetroArch/retroarch "$APPDIR/"; \
  chmod +x "$APPDIR/retroarch"; \
  cp -r /src/retroarch-assets/* "$APPDIR/assets/"; \
  cp -r /src/retroarch-joypad-autoconfig/* "$APPDIR/autoconfig/"; \
  cp /src/cores-built/*.so "$APPDIR/cores/" || true; \
  # Prepare to wrap the ELF later
  mv "$APPDIR/retroarch" "$APPDIR/retroarch.bin"

# Build SDL2 for aarch64 with fbcon + kmsdrm (PulseAudio disabled) and bundle it
RUN set -eux; \
  SDL_VER=2.30.7; \
  cd /src; \
  curl -L -o SDL2.tar.gz "https://github.com/libsdl-org/SDL/releases/download/release-$SDL_VER/SDL2-$SDL_VER.tar.gz"; \
  tar -xzf SDL2.tar.gz; \
  cd "SDL2-$SDL_VER"; \
  PKG_CONFIG_LIBDIR=/usr/lib/aarch64-linux-gnu/pkgconfig \
  ./configure --host="$TARGET" --prefix=/opt/sdl2-aarch64 \
    --disable-rpath \
    --enable-video-kmsdrm \
    --enable-video-fbcon \
    --disable-video-wayland \
    --disable-video-x11 \
    --disable-pulseaudio \
    --enable-alsa; \
  make -j"$(nproc)" CC="$TARGET-gcc"; \
  make install; \
  # Bundle our SDL2 runtime (prefer it over distro one)
  cp -Lf /opt/sdl2-aarch64/lib/libSDL2-2.0.so.0 "$APPDIR/lib/"

# Trim bulky asset dirs
RUN set -eux; \
  rm -rf "$APPDIR/assets/wallpapers" || true; \
  rm -rf "$APPDIR/assets/xmb" || true; \
  rm -rf "$APPDIR/assets/src" || true

# Bundle dynamic linker and essential runtime libs + some common deps
RUN set -eux; \
  APP_LIB="$APPDIR/lib"; mkdir -p "$APP_LIB" "$APPDIR/usr/lib"; \
  LIBDIRS="/lib/aarch64-linux-gnu /usr/lib/aarch64-linux-gnu /usr/aarch64-linux-gnu/lib"; \
  # Dynamic linker
  DL="$($TARGET-gcc -print-file-name=ld-linux-aarch64.so.1)"; \
  [ -f "$DL" ] || DL="$(find $LIBDIRS -name 'ld-linux-aarch64.so.1' -type f -print -quit)"; \
  [ -f "$DL" ]; \
  cp "$DL" "$APP_LIB/"; \
  DL_DIR="$(dirname "$(readlink -f "$DL")")"; \
  # glibc core from same dir as the loader to keep versions in sync
  for L in libc.so.6 libm.so.6 libpthread.so.0 librt.so.1 libdl.so.2; do \
    if [ -f "$DL_DIR/$L" ]; then cp "$DL_DIR/$L" "$APP_LIB/"; \
    else find $LIBDIRS -maxdepth 1 -name "$L" -type f -print -quit | xargs -r -I{} cp "{}" "$APP_LIB/"; \
    fi; \
  done; \
  # C++/GCC runtimes
  cp "$($TARGET-g++ -print-file-name=libstdc++.so.6)" "$APP_LIB/"; \
  cp "$($TARGET-gcc -print-file-name=libgcc_s.so.1)" "$APP_LIB/"; \
  # Seed likely runtime libs (the rest are auto-copied below)
  # NOTE: we intentionally do NOT seed libSDL2-2.0.so.0 here, as we bundled our custom SDL2 above.
  for L in \
    libasound.so.2 libudev.so.1 \
    libEGL.so.1 libGLESv2.so.2 libGL.so.1 \
    libvulkan.so.1 libgbm.so.1 libdrm.so.2 \
    libz.so.1 libpng16.so.16 libjpeg.so.62 libfreetype.so.6 \
    libssl.so.3 \
  ; do \
    for d in /lib/aarch64-linux-gnu /usr/lib/aarch64-linux-gnu /lib /usr/lib; do \
      if [ -e "$d/$L" ]; then cp -Lf "$d/$L" "$APP_LIB/" && break; fi; \
    done; \
  done; \
  # PulseAudio bits (in case any path ends up requiring libpulse at runtime)
  for d in /lib/aarch64-linux-gnu /usr/lib/aarch64-linux-gnu; do \
    [ -e "$APP_LIB/libpulse.so.0" ] || { [ -e "$d/libpulse.so.0" ] && cp -Lf "$d/libpulse.so.0" "$APP_LIB/"; }; \
    [ -e "$APP_LIB/libpulse-simple.so.0" ] || { [ -e "$d/libpulse-simple.so.0" ] && cp -Lf "$d/libpulse-simple.so.0" "$APP_LIB/"; }; \
  done; \
  copied_common=0; \
  for p in /usr/lib/aarch64-linux-gnu/pulseaudio/libpulsecommon-*.so /lib/aarch64-linux-gnu/pulseaudio/libpulsecommon-*.so; do \
    if [ -e "$p" ]; then cp -Lf "$p" "$APP_LIB/"; copied_common=1; break; fi; \
  done; \
  if [ -e "$APP_LIB/libpulse.so.0" ] && [ "$copied_common" -eq 0 ]; then \
    echo "WARN: libpulse.so.0 present but libpulsecommon-*.so not found; runtime may fail."; \
  fi; \
  # Verify presence of critical libs (fail fast if missing)
  for MUST in ld-linux-aarch64.so.1 libc.so.6 libm.so.6 libpthread.so.0 libdl.so.2 libstdc++.so.6 libgcc_s.so.1 libasound.so.2 libSDL2-2.0.so.0; do \
    test -f "$APP_LIB/$MUST" || { echo "ERROR: missing $MUST in $APP_LIB"; exit 1; }; \
  done

# Auto-copy runtime deps via readelf (cross-safe) and verify none are missing
RUN set -eux; \
  APPDIR="/opt/appdir"; APP_LIB="$APPDIR/lib"; mkdir -p "$APP_LIB"; \
  LIB_DIRS="/lib/aarch64-linux-gnu /usr/lib/aarch64-linux-gnu /usr/aarch64-linux-gnu/lib"; \
  neededs() { \
    readelf -d "$1" 2>/dev/null | awk '/NEEDED/ { gsub(/[\[\]]/,"",$5); print $5 }' | sort -u; \
  }; \
  copy_needed_by_name() { \
    local name="$1"; \
    [ -e "$APP_LIB/$name" ] && return 0; \
    local path=""; \
    for d in $LIB_DIRS; do \
      [ -e "$d/$name" ] && { path="$d/$name"; break; } || true; \
    done; \
    if [ -n "$path" ]; then cp -Lf "$path" "$APP_LIB/$name"; return 0; fi; \
    if command -v ldconfig >/dev/null 2>&1; then \
      p="$(ldconfig -p 2>/dev/null | awk -v n="$name" '$1==n {print $NF; exit}')" || true; \
      [ -n "$p" ] && [ -e "$p" ] && { cp -Lf "$p" "$APP_LIB/$name"; return 0; }; \
    fi; \
    echo "WARN: could not locate $name"; return 1; \
  }; \
  process_file() { \
    local f="$1"; [ -f "$f" ] || return 0; \
    neededs "$f" | while read -r n; do [ -n "$n" ] && copy_needed_by_name "$n" || true; done; \
  }; \
  targets="$APPDIR/retroarch.bin"; \
  if [ -d "$APPDIR/cores" ]; then \
    cores="$(find "$APPDIR/cores" -maxdepth 1 -type f -name '*.so' -print || true)"; \
    [ -n "$cores" ] && targets="$targets $cores"; \
  fi; \
  for round in 1 2 3 4; do \
    changed=0; before="$(ls -1 "$APP_LIB" 2>/dev/null | wc -l || echo 0)"; \
    for t in $targets $(find "$APP_LIB" -maxdepth 1 -type f -name 'lib*.so*' -print 2>/dev/null || true); do \
      process_file "$t"; \
    done; \
    after="$(ls -1 "$APP_LIB" 2>/dev/null | wc -l || echo 0)"; \
    [ "$after" -gt "$before" ] && changed=1; \
    [ "$changed" -eq 0 ] && break; \
  done; \
  missing=0; \
  for n in $(neededs "$APPDIR/retroarch.bin"); do \
    [ -e "$APP_LIB/$n" ] || { echo "ERROR: unresolved dependency for retroarch.bin: $n"; missing=1; }; \
  done; \
  [ "$missing" -eq 0 ] || { ls -lah "$APP_LIB" || true; exit 1; }; \
  echo "Dependency copy and verification completed."

# Default config to use SDL2 video/input and ALSA audio, and disable threaded video
RUN set -eux; \
  cat > "$APPDIR/retroarch.cfg" <<'CFG'
menu_driver = "ozone"

video_driver = "sdl2"
video_threaded = "false"
audio_driver = "alsa"
input_driver = "udev"

assets_directory = "assets"
config_save_on_exit = "true"
core_directory = "cores"
libretro_directory = "cores"
rgui_browser_directory = "downloads"
joypad_autoconfig_dir = "autoconfig"
playlist_directory = "playlists"
recording_output_directory = "records"
shader_directory = "shaders"
system_directory = "system"

video_fullscreen = "true"

input_enable_hotkey = "escape"
input_exit_emulator = "q"

udev_enable = "true"
CFG

# Wrapper launcher using bundled dynamic linker and robust environment setup
RUN set -eux; \
  cat > "$APPDIR/retroarch" <<'SH'
#!/bin/sh
SELF_DIR="$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)"
unset LD_PRELOAD LD_AUDIT

# Ensure XDG_RUNTIME_DIR is valid
if [ -z "${XDG_RUNTIME_DIR:-}" ]; then
  UID_DIR="/run/user/$(id -u)"
  if [ -d "$UID_DIR" ]; then
    export XDG_RUNTIME_DIR="$UID_DIR"
  else
    export XDG_RUNTIME_DIR="$SELF_DIR/.xdg-runtime"
    mkdir -p "$XDG_RUNTIME_DIR"
    chmod 700 "$XDG_RUNTIME_DIR" 2>/dev/null || true
  fi
fi

# Choose SDL video driver:
# If a desktop is present, let SDL decide (x11/wayland).
# Otherwise prefer fbcon if /dev/fb0 exists; else kmsdrm if /dev/dri exists; else dummy.
if [ -z "${SDL_VIDEODRIVER:-}" ]; then
  if [ -n "${DISPLAY:-}" ] || [ -n "${WAYLAND_DISPLAY:-}" ]; then
    : # desktop present
  elif [ -e /dev/fb0 ]; then
    export SDL_VIDEODRIVER=fbcon
    export SDL_FBDEV=/dev/fb0
  elif ls /dev/dri/card* >/dev/null 2>&1; then
    export SDL_VIDEODRIVER=kmsdrm
  else
    export SDL_VIDEODRIVER=dummy
  fi
fi

# Prefer ALSA to avoid Pulse issues
export SDL_AUDIODRIVER="${SDL_AUDIODRIVER:-alsa}"

export LD_LIBRARY_PATH="$SELF_DIR/lib:$SELF_DIR/usr/lib"
exec "$SELF_DIR/lib/ld-linux-aarch64.so.1" \
  --library-path "$LD_LIBRARY_PATH" \
  "$SELF_DIR/retroarch.bin" --video-threaded=false "$@"
SH
RUN chmod +x "$APPDIR/retroarch"

# Also provide a run.sh with the same behavior (optional convenience)
RUN set -eux; \
  cat > "$APPDIR/run.sh" <<'SH'
#!/bin/sh
SELF_DIR="$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)"
unset LD_PRELOAD LD_AUDIT

# Ensure XDG_RUNTIME_DIR is valid
if [ -z "${XDG_RUNTIME_DIR:-}" ]; then
  UID_DIR="/run/user/$(id -u)"
  if [ -d "$UID_DIR" ]; then
    export XDG_RUNTIME_DIR="$UID_DIR"
  else
    export XDG_RUNTIME_DIR="$SELF_DIR/.xdg-runtime"
    mkdir -p "$XDG_RUNTIME_DIR"
    chmod 700 "$XDG_RUNTIME_DIR" 2>/dev/null || true
  fi
fi

# Choose SDL video driver with the same logic
if [ -z "${SDL_VIDEODRIVER:-}" ]; then
  if [ -n "${DISPLAY:-}" ] || [ -n "${WAYLAND_DISPLAY:-}" ]; then
    : # desktop present
  elif [ -e /dev/fb0 ]; then
    export SDL_VIDEODRIVER=fbcon
    export SDL_FBDEV=/dev/fb0
  elif ls /dev/dri/card* >/dev/null 2>&1; then
    export SDL_VIDEODRIVER=kmsdrm
  else
    export SDL_VIDEODRIVER=dummy
  fi
fi

export SDL_AUDIODRIVER="${SDL_AUDIODRIVER:-alsa}"
export LD_LIBRARY_PATH="$SELF_DIR/lib:$SELF_DIR/usr/lib"
exec "$SELF_DIR/lib/ld-linux-aarch64.so.1" \
  --library-path "$LD_LIBRARY_PATH" \
  "$SELF_DIR/retroarch.bin" --video-threaded=false "$@"
SH
RUN chmod +x "$APPDIR/run.sh"

# Package result
RUN set -eux; \
  TAR_NAME=retroarch-aarch64-tinalinux-portable.tar.gz; \
  pushd "$(dirname "$APPDIR")"; \
    tar -czf "$OUTDIR/$TAR_NAME" "$(basename "$APPDIR")"; \
  popd; \
  echo "Created $OUTDIR/$TAR_NAME"

CMD ["/bin/bash", "-lc", "ls -lah /out && echo 'Use: docker create && docker cp to extract the tarball.'"]
