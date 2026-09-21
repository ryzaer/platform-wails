import sys
import yaml


def get_serial_file():
    try:
        index = sys.argv.index("--serial")

        if index + 1 < len(sys.argv):
            return sys.argv[index + 1]

    except ValueError:
        pass

    return None


SERIAL_FILE = get_serial_file()


def output(*args, **kwargs):
    text = " ".join(
        str(arg)
        for arg in args
    )

    end = kwargs.get("end", "\n")

    if SERIAL_FILE:
        with open(
            SERIAL_FILE,
            "a",
            encoding="utf-8"
        ) as f:
            f.write(text + end)

    else:
        print(
            *args,
            **kwargs
        )


def hello():
    output("Hello from Python renderer!")


def render_invoice(yaml_file, output_file):
    with open(
        yaml_file,
        "r",
        encoding="utf-8"
    ) as f:
        data = yaml.safe_load(f)

    invoice = data.get("invoice", {})

    output("=== INVOICE ===")
    output("Invoice :", invoice.get("number"))
    output("Customer:", invoice.get("customer"))
    output("Total   :", invoice.get("total"))

    with open(
        output_file,
        "w",
        encoding="utf-8"
    ) as f:
        f.write("INVOICE\n")
        f.write(
            f"Number: {invoice.get('number')}\n"
        )
        f.write(
            f"Customer: {invoice.get('customer')}\n"
        )
        f.write(
            f"Total: {invoice.get('total')}\n"
        )


def generate_qr(data):
    output("=== QR ===")
    output("Data:", data)


def render_pdf(filename):
    output("=== PDF ===")
    output("File:", filename)


def main():
    output("Python renderer started")
    output("ARGV:", sys.argv)

    if len(sys.argv) < 2:
        hello()
        return 0

    command = sys.argv[1]

    if command == "--core":

        if len(sys.argv) < 3:
            print("Missing core command")
            return 1

        action = sys.argv[2]

        if action == "print":

            if len(sys.argv) < 4:
                print("Usage: --core print <yaml> [output]")
                return 1

            yaml_file = sys.argv[3]

            if len(sys.argv) >= 5:
                output_file = sys.argv[4]
            else:
                output_file = "invoice.txt"

            render_invoice(
                yaml_file,
                output_file
            )

            return 0

        if action == "qr":

            if len(sys.argv) < 4:
                print("Usage: --core qr <data>")
                return 1

            generate_qr(
                sys.argv[3]
            )

            return 0

        if action == "pdf":

            if len(sys.argv) < 4:
                print("Usage: --core pdf <file>")
                return 1

            render_pdf(
                sys.argv[3]
            )

            return 0

        if action == "hello":
            hello()
            return 0

        print(
            f"Unknown core command: {action}"
        )

        return 1

    print(
        f"Unknown command: {command}"
    )

    return 1


if __name__ == "__main__":
    raise SystemExit(main())