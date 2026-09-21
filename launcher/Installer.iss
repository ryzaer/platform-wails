#define MyAppName "Elsana"
#define MyAppVersion "1.0.0"
#define MyAppPublisher "Elsana"
#define MyAppExeName "elsana.exe"

[Setup]
AppId={{7D2B4F9A-1E63-4F2A-9B71-5C8E3A6D2041}
AppName={#MyAppName}
AppVersion={#MyAppVersion}
AppPublisher={#MyAppPublisher}

DefaultDirName={autopf}\Elsana
DefaultGroupName=Elsana

SetupIconFile=Elsana.ico

OutputDir=..\integrator\build\bin
OutputBaseFilename=Elsana-Setup

ArchitecturesAllowed=x64compatible
ArchitecturesInstallIn64BitMode=x64compatible

PrivilegesRequired=admin

Compression=lzma
SolidCompression=yes

WizardStyle=modern
WizardImageFile=164x314.bmp
WizardSmallImageFile=55x55.bmp

UninstallDisplayIcon={app}\elsana.exe

[Languages]
Name: "indonesian"; MessagesFile: "compiler:Languages\Indonesian.isl"

[Files]
Source: "..\integrator\build\bin\windows\*"; \
    DestDir: "{app}"; \
    Flags: ignoreversion recursesubdirs createallsubdirs

[Icons]
Name: "{group}\Elsana"; \
    Filename: "{app}\elsana.exe"; \
    Tasks: startmenu

Name: "{autodesktop}\Elsana"; \
    Filename: "{app}\elsana.exe"; \
    Tasks: desktop

[Tasks]
Name: "startmenu"; \
    Description: "Buat shortcut di Start Menu"; \
    GroupDescription: "Shortcuts:"

Name: "desktop"; \
    Description: "Buat shortcut di Desktop"; \
    GroupDescription: "Shortcuts:"

[Run]
Filename: "{app}\elsana.exe"; \
    Description: "Launch Elsana"; \
    Flags: nowait postinstall skipifsilent