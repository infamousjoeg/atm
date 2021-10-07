# atm
An ATM in the console

## Usage
Start the ./atm binary and then enter in one of the following commands
```
authorize <account_id> <pin>
withdraw <amount>
deposit <amount>
balance
history
logout
end
```

## Development
To compile the code execute execute the following command in the root directory:
```bash
go build -o atm cmd/atm/main.go
```

A file called `./atm` will appear in your current working directory to run this file simply call it:
```bash
./atm
```

## Testing
To run all tests in the suite execute the following in the root directory:
```bash
go test
```
