#include <windows.h>
#include <objbase.h>
#include <shellapi.h>

#include <iostream>
#include <filesystem>
#include <string>
#include <vector>

using namespace std;
namespace fs = std::filesystem;


// ============================================================
// DLL FUNCTION TYPES
// ============================================================

using StartIntegratorFunc = void (*)();

using RunCodeFunc =
    int (__stdcall *)(
        int,
        wchar_t**,
        const wchar_t*
    );


// ============================================================
// GLOBAL
// ============================================================


// ============================================================
// GET APPLICATION DIRECTORY
// ============================================================

fs::path GetAppDirectory()
{
    wchar_t buffer[MAX_PATH];

    DWORD length = GetModuleFileNameW(
        nullptr,
        buffer,
        MAX_PATH
    );

    if (length == 0 || length >= MAX_PATH)
    {
        return {};
    }

    return fs::path(buffer).parent_path();
}

// ============================================================
// LOAD DLL
// ============================================================

HMODULE LoadModule(
    const fs::path& path,
    bool silent = false
)
{
    if (!silent)
    {
        std::wcout
            << L"Loading: "
            << path
            << L"\n";
    }

    HMODULE dll =
        LoadLibraryW(path.c_str());

    if (!dll)
    {
        DWORD error = GetLastError();

        std::cerr
            << "Failed to load DLL. Error: "
            << error
            << "\n";

        return nullptr;
    }

    if (!silent)
    {
        std::cout
            << "DLL loaded.\n";
    }

    return dll;
}


// ============================================================
// RUN PYTHON / NUITKA
// ============================================================
void AttachParentConsole()
{
    if (!AttachConsole(ATTACH_PARENT_PROCESS))
    {
        return;
    }

    FILE* fp = nullptr;

    freopen_s(
        &fp,
        "CONOUT$",
        "w",
        stdout
    );

    freopen_s(
        &fp,
        "CONOUT$",
        "w",
        stderr
    );

    freopen_s(
        &fp,
        "CONIN$",
        "r",
        stdin
    );

    std::ios::sync_with_stdio(true);
}

int RunRenderer(
    const vector<wstring>& args
)
{
    fs::path dllPath =
        GetAppDirectory() /
        L"renderer.dll";

    vector<wstring> rendererArgs =
        args;

    if (rendererArgs.size() < 2)
    {
        return 1;
    }

    bool silent = false;

    for (size_t i = 2; i < rendererArgs.size(); ++i)
    {
        if (rendererArgs[i] == L"--serial")
        {
            silent = true;
            break;
        }
    }

    rendererArgs[0] =
        dllPath.wstring();

    HMODULE dll =
        LoadModule(
            dllPath,
            silent
        );

    if (!dll)
    {
        return 1;
    }

    auto run_code =
        reinterpret_cast<RunCodeFunc>(
            GetProcAddress(
                dll,
                "run_code"
            )
        );

    if (!run_code)
    {
        FreeLibrary(dll);
        return 1;
    }

    if (!silent)
    {
        std::cout
            << "run_code found.\n";
    }

    vector<wchar_t*> argv;

    for (auto& arg : rendererArgs)
    {
        argv.push_back(
            const_cast<wchar_t*>(
                arg.c_str()
            )
        );
    }

    if (!silent)
    {
        std::cout
            << "Calling run_code...\n";

        std::cout.flush();
    }

    int result =
        run_code(
            static_cast<int>(
                argv.size()
            ),
            argv.data(),
            dllPath.c_str()
        );

    if (!silent)
    {
        std::cout
            << "run_code returned: "
            << result
            << "\n";

        std::cout.flush();
    }

    FreeLibrary(dll);

    return result;
}

// ============================================================
// INTEGRATOR
// ============================================================

int RunIntegrator()
{
    fs::path dllPath =
        GetAppDirectory() /
        L"integrator.dll";

    HMODULE dll = LoadModule(dllPath);

    if (!dll)
    {
        return 1;
    }

    auto StartIntegrator =
        reinterpret_cast<StartIntegratorFunc>(
            GetProcAddress(
                dll,
                "StartIntegrator"
            )
        );

    if (!StartIntegrator)
    {
        DWORD error = GetLastError();

        std::cerr
            << "StartIntegrator not found. Error: "
            << error
            << "\n";

        FreeLibrary(dll);

        return 1;
    }

    std::cout
        << "StartIntegrator found.\n";

    StartIntegrator();

    FreeLibrary(dll);

    return 0;
}


// ============================================================
// COMMAND LINE
// ============================================================

vector<wstring> GetArguments(
    int argc,
    wchar_t* argv[]
)
{
    vector<wstring> args;

    for (int i = 0; i < argc; ++i)
    {
        args.emplace_back(argv[i]);
    }

    return args;
}


// ============================================================
// PRINT HELP
// ============================================================

void PrintHelp()
{
    std::wcout << L"\n";
    std::wcout << L"Elsana Launcher\n";
    std::wcout << L"\n";

    std::wcout << L"Usage:\n";
    std::wcout << L"  launcher.exe\n";
    std::wcout << L"  launcher.exe --gui\n";
    std::wcout << L"  launcher.exe --core <command> [arguments]\n";
    std::wcout << L"  launcher.exe --core <command> [arguments] [--serial <file>]\n";

    std::wcout << L"\n";

    std::wcout << L"Examples:\n";
    std::wcout << L"  launcher.exe --core hello\n";
    std::wcout << L"  launcher.exe --core hello --serial tmp\\serial-XXXXXXX.txt\n";
    std::wcout << L"  launcher.exe --core print invoice.yml\n";
    std::wcout << L"  launcher.exe --core print invoice.yml --serial tmp\\serial-XXXXXXX.txt\n";
    std::wcout << L"  launcher.exe --core qr \"https://elsana.app\"\n";
    std::wcout << L"  launcher.exe --core pdf invoice.yml\n";

    std::wcout << L"\n";
}


// ============================================================
// MAIN
// ============================================================

int WINAPI wWinMain(
    HINSTANCE,
    HINSTANCE,
    PWSTR,
    int
)
{
    int argc;

    LPWSTR* argv =
        CommandLineToArgvW(
            GetCommandLineW(),
            &argc
        );

    if (!argv)
    {
        return 1;
    }


    // --------------------------------------------------------
    // COM
    // --------------------------------------------------------

    HRESULT hr = CoInitializeEx(
        nullptr,
        COINIT_APARTMENTTHREADED
    );

    if (FAILED(hr))
    {
        std::cerr
            << "CoInitializeEx failed. HRESULT: 0x"
            << std::hex
            << hr
            << std::dec
            << "\n";

        LocalFree(argv);

        return 1;
    }


    // --------------------------------------------------------
    // ARGUMENTS
    // --------------------------------------------------------

    auto args = GetArguments(
        argc,
        argv
    );


    // --------------------------------------------------------
    // DEFAULT
    // --------------------------------------------------------

    if (argc == 1)
    {
        int result = RunIntegrator();

        CoUninitialize();
        LocalFree(argv);

        return result;
    }


    // --------------------------------------------------------
    // COMMAND
    // --------------------------------------------------------

    const wstring& command = args[1];

    // --------------------------------------------------------
    // ATTACH CONSOLE FOR GUI / HELP
    // --------------------------------------------------------

    if (
        command == L"--gui" ||
        command == L"--help" ||
        command == L"-h" ||
        command == L"/?"  ||
        command == L"--core"
    )
    {
        if (AttachConsole(ATTACH_PARENT_PROCESS))
        {
            FILE* fp = nullptr;

            freopen_s(
                &fp,
                "CONOUT$",
                "w",
                stdout
            );

            freopen_s(
                &fp,
                "CONOUT$",
                "w",
                stderr
            );

            freopen_s(
                &fp,
                "CONIN$",
                "r",
                stdin
            );

            std::ios::sync_with_stdio(true);
        }
    }


    // --------------------------------------------------------
    // GUI
    // --------------------------------------------------------

    if (command == L"--gui")
    {
        int result = RunIntegrator();

        CoUninitialize();
        LocalFree(argv);

        return result;
    }


    // --------------------------------------------------------
    // CORE
    // --------------------------------------------------------
    if (command == L"--core")
    {
        int result = RunRenderer(args);

        CoUninitialize();
        LocalFree(argv);

        return result;
    }


    // --------------------------------------------------------
    // HELP
    // --------------------------------------------------------

    if (
        command == L"--help" ||
        command == L"-h" ||
        command == L"/?"
    )
    {
        PrintHelp();

        CoUninitialize();
        LocalFree(argv);

        return 0;
    }


    // --------------------------------------------------------
    // UNKNOWN COMMAND
    // --------------------------------------------------------

    std::wcerr
        << L"Unknown command: "
        << command
        << L"\n";

    PrintHelp();

    CoUninitialize();
    LocalFree(argv);

    return 0;
}