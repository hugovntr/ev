#!/usr/bin/env bash
set -e

# Configuration
APP_NAME="ev"
VERSION=${1:-"v1.0.0"} # Version from argument
BUILD_DIR="dist"
PLATFORMS=(
  "linux/amd64"
  "linux/arm64"
  "darwin/amd64"
  "darwin/arm64"
  "windows/amd64"
)

echo "🚀 Starting release build for version: $VERSION"

# Clean build directory
rm -rf $BUILD_DIR
mkdir -p $BUILD_DIR

for PLATFORM in "${PLATFORMS[@]}"; do
  IFS='/' read -r -a PARTS <<<"$PLATFORM"
  GOOS=${PARTS[0]}
  GOARCH=${PARTS[1]}
  OUTPUT_NAME=$APP_NAME
  if [ "$GOOS" == "windows" ]; then
    OUTPUT_NAME+=".exe"
  fi
  echo "   🔨 Building for $GOOS/$GOARCH..."
  env GOOS=$GOOS GOARCH=$GOARCH go build -ldflags="-s -w" -o $OUTPUT_NAME main.go
  ARCHIVE_NAME="${APP_NAME}_${VERSION}_${GOOS}_${GOARCH}"
  if [ "$GOOS" == "windows" ]; then
    zip -q "$BUILD_DIR/${ARCHIVE_NAME}.zip" $OUTPUT_NAME
  else
    tar -czf "$BUILD_DIR/${ARCHIVE_NAME}.tar.gz" $OUTPUT_NAME
  fi
  rm $OUTPUT_NAME
done

echo "   📝 Generating checksums..."
cd $BUILD_DIR
shasum -a 256 * >checksums.txt

echo "✅ Release ready in '$BUILD_DIR' directory!"
echo "------------------------------------------------"
ls -lh
