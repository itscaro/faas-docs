build:
	docker run --rm -ti --user root -v "$(shell pwd):/docs" -w /docs -p 8000:8000 itscaro/mkdocs build

serve:
	docker run --rm -ti --user root -v "$(shell pwd):/docs" -w /docs -p 8000:8000 itscaro/mkdocs serve

