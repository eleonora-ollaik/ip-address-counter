# IPv4 Addresses counter
Golang program that can count and output unique number of IPv4 addresses from a given txt file.

## Installation

Clone this repository:

```bash
git clone https://github.com/eleonora-ollaik/ip-address-counter.git
```

## Usage

When using you **must** provide a filename when you run the programm. Run the following command from the main directory: 

```bash
go run . <filename>
```

For example: 

```bash
go run . ip_addresses.txt
```

## Testing:

For a smaller subset use *ip_addresses.txt*.
It has roughly ~1500 ip addresses from which only 79 are unique.
If you want to test it with the large ip_addresses file, you need to put it unzipped in the same directory first.

### Expected output:

#### Happy scenario:

```bash
go run . ip_addresses.txt
```

    Starting the count...


    Execution complete.

    Unique number of IPs: 79
    Lines processed: 1591
    Time taken to execute: 2.0375ms


```bash
go run . ip_addresses
```


    Starting the count...

    Progress: 99.86% (108988 /109139 MB), 104.7 MB/s
    Execution complete.

    Unique number of IPs: 1000000000
    Lines processed: 7999999994
    Time taken to execute: 17m21.702455s


#### Unhappy scenarios:

```bash
go run . non-existent-file
```

    Error: file ip_addresses does not exist

-----------------

The following file contains 3 strings instead of IPv4 addresses
The program is designed to log the error and skip those lines to continue count

```bash
go run . ip_addresses_invalids.txt
```


    Starting the count...

    bebebe is not a valid IPv4 address
    lalala is not a valid IPv4 address
    hahaha is not a valid IPv4 address

    Execution complete.

    Unique number of IPs: 79
    Lines processed: 1592
    Time taken to execute: 2.5909ms
