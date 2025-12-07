build: 
	go build -o advent-of-code-2025 main.go

runall: 
	./advent-of-code-2025 run -d 1 -i inputs/day1.txt -q
	./advent-of-code-2025 run -d 2 -i inputs/day2.txt -q
	./advent-of-code-2025 run -d 3 -i inputs/day3.txt -q
	./advent-of-code-2025 run -d 4 -i inputs/day4.txt -q
	./advent-of-code-2025 run -d 5 -i inputs/day5.txt -q
	./advent-of-code-2025 run -d 6 -i inputs/day6.txt -q
	./advent-of-code-2025 run -d 7 -i inputs/day7.txt -q
