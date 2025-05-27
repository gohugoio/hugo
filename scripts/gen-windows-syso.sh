#!/usr/bin/env bash
set -e

# Determine version from git tags or fallback to "0.0.0"
VERSION=$(git describe --tags --abbrev=0 2>/dev/null || echo "0.0.0")
GOVERSIONINFO=${GOVERSIONINFO:-goversioninfo}

cat > versioninfo.json <<EOF
{
  "FixedFileInfo": {
    "FileVersion": "${VERSION}.0",
    "ProductVersion": "${VERSION}.0",
    "FileFlagsMask": "3f",
    "FileFlags": "00",
    "FileOS": "040004",
    "FileType": "01",
    "FileSubtype": "00"
  },
  "StringFileInfo": {
    "Comments": "",
    "CompanyName": "gohugoio",
    "FileDescription": "Hugo - The world’s fastest framework for building websites",
    "FileVersion": "${VERSION}.0",
    "InternalName": "hugo",
    "LegalCopyright": "Copyright $(date +%Y) The Hugo Authors",
    "OriginalFilename": "hugo.exe",
    "ProductName": "Hugo",
    "ProductVersion": "${VERSION}.0"
  },
  "VarFileInfo": {
    "Translation": {
      "LangID": "0409",
      "CharsetID": "04B0"
    }
  },
  "IconPath": ""
}
EOF

# Ensure goversioninfo is installed
if ! command -v $GOVERSIONINFO &>/dev/null; then
  echo "goversioninfo not found. Please install it: go install github.com/josephspurrier/goversioninfo/cmd/goversioninfo@latest"
  exit 1
fi

$GOVERSIONINFO -product-version "${VERSION}" -file-version "${VERSION}" -o resource_windows.syso