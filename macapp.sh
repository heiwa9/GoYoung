#!/bin/sh

APP="GoYoung.app"
mkdir -p $APP/Contents/{MacOS,Resources}
go build -o $APP/Contents/MacOS/goyoung
cat > $APP/Contents/Info.plist << EOF
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
	<key>CFBundleExecutable</key><string>goyoung</string>
	<key>CFBundleGetInfoString</key><string>share files internally</string>
	<key>CFBundleIconFile</key><string>myapp</string>
	<key>CFBundleIdentifier</key><string>sharefile.geekland.cc</string>
	<key>CFBundlePackageType</key><string>APPL</string>
	<key>CFBundleSignature</key><string>????</string>
	<key>LSMinimumSystemVersion</key><string>10.13</string>
	<key>NOTE</key><string>share files internally</string>
	<key>NSPrincipalClass</key><string>NSApplication</string>
	<key>LSUIElement</key><true/><key>NSSupportsAutomaticGraphicsSwitching</key><true/>
</dict>
</plist>
EOF
cp icon/myapp.png $APP/Contents/Resources/myapp.png
find $APP
