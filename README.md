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
go run main.go <filename>
```

For example: 

```bash
go run main.go ip_addresses.txt
```

## Testing:

For a smaller subset use *ip_addresses.txt*.
It has roughly ~1400 ip addresses from which only 79 are unique.

### Expected output:

#### Happy scenario:

```bash
go run main.go ip_addresses.txt
```


    Starting the count...


    Execution complete.


    Unique number of IPs: 79

    Lines processed: 1472

    Time taken to execute: 4.3383ms



#### Unhappy scenarios:

```bash
go run main.go non-existent-file
```

    Error: file ip_addresses does not exist

-----------------

The following file contains 3 strings instead of IPv4 addresses
The program is designed to log the error and skip those lines to continue count

```bash
go run main.go ip_addresses_invalids.txt
```


    Starting the count...


    hahaha is not a valid IPv4 address
    lalala is not a valid IPv4 address
    bebebe is not a valid IPv4 address


    Execution complete.


    Unique number of IPs: 79

    Lines processed: 1472

    Time taken to execute: 4.2913ms
