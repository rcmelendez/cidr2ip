# cidr2ip

`cidr2ip` is a lightweight CLI tool that transforms CIDR notations into a list of IP addresses. The generated list is saved to a CSV file, making it convenient for creating [lookup tables](https://docs.devo.com/space/latest/95204003/Upload+a+lookup+table) for Devo.

## Benefits

- Incredibly easy to use.
- Generate a list of IP addresses from multiple CIDR notations.
- Support for both command-line arguments and input from a file.
- Includes all IP addresses (network and broadcast addresses included).

## Getting Started

### Installation

`cidr2ip` is available on Linux, macOS, and Windows platforms. The pre-built binaries are available in the [Releases](https://github.com/rcmelendez/cidr2ip/releases) page.

Alternatively, you can build it from source by running:

```bash
go build .
```

### Usage

For macOS/Linux:
```bash
./cidr2ip [-f filename] <CIDR1 CIDR2 ...>
```

For Windows:
```shell
.\cidr2ip.exe [-f filename] <CIDR1 CIDR2 ...>
```
> Pro tip: To avoid using the full path, add the directory containing the `cidr2ip` binary to the system's PATH environment variable.

## Examples

Generate IP list from CIDRs `192.168.1.0/24` and `10.0.0.0/8` as command-line arguments:
```bash
cidr2ip 192.168.1.0/24 10.0.0.0/8
```

Generate IP list from the file `cidr_list`:
```bash
cidr2ip -f cidr_list
```

## Successful Output

Upon successful execution, `cidr2ip` will display a message indicating the generated CSV filename along with a timestamp. 
The format is as follows:
```plaintext
IP list saved to cidr2ip_YYYY-MM-DD_HH-MM-SS.csv
```

## License

`cidr2ip` is licensed under the terms of the [MIT License](https://github.com/rcmelendez/cidr2ip/blob/main/LICENSE).

## Contact

Find me as __rcmelendez__ on [LinkedIn](https://www.linkedin.com/in/rcmelendez/), [Medium](https://rcmelendez.medium.com/), and [GitHub](https://github.com/rcmelendez/).