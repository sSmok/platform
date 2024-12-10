package db

//go:generate sh -c "rm -rf mocks && mkdir -p mocks"
//go:generate minimock -i TransactorI -o ./mocks/ -s "_minimock.go"
