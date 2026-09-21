#define MyAppName "Integrator"
#define MyAppVersion "1.0.0"
#define MyAppPublisher "Integrator"
#define MyAppExeName "app.exe"

[Setup]
AppId={{7D2B4F9A-1E63-4F2A-9B71-5C8E3A6D2041}
AppName={#MyAppName}
AppVersion={#MyAppVersion}
AppPublisher={#MyAppPublisher}

DefaultDirName={autopf}\{#MyAppName}
DefaultGroupName={#MyAppPublisher}

SetupIconFile=Installer.ico

OutputDir=..\integrator\build\bin
OutputBaseFilename=app-setup

ArchitecturesAllowed=x64compatible
ArchitecturesInstallIn64BitMode=x64compatible

PrivilegesRequired=admin

Compression=lzma
SolidCompression=yes

WizardStyle=modern
WizardImageFile=InstallerBg-164x314.bmp
WizardSmallImageFile=InstallerBg-55x55.bmp

UninstallDisplayIcon={app}\{#MyAppName}

[Languages]
Name: "indonesian"; MessagesFile: "compiler:Languages\Indonesian.isl"

[Files]
Source: "..\integrator\build\bin\windows\*"; \
    DestDir: "{app}"; \
    Flags: ignoreversion recursesubdirs createallsubdirs

[Icons]
Name: "{group}\{#MyAppPublisher}"; \
    Filename: "{app}\{#MyAppName}"; \
    Tasks: startmenu

Name: "{autodesktop}\{#MyAppPublisher}"; \
    Filename: "{app}\{#MyAppName}"; \
    Tasks: desktop

[Tasks]
Name: "startmenu"; \
    Description: "Buat shortcut di Start Menu"; \
    GroupDescription: "Shortcuts:"

Name: "desktop"; \
    Description: "Buat shortcut di Desktop"; \
    GroupDescription: "Shortcuts:"

[Run]
Filename: "{app}\{#MyAppName}"; \
    Description: "Launch {#MyAppName}"; \
    Flags: nowait postinstall skipifsilent