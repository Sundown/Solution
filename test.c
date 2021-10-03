#include <stdio.h>
#include <stdlib.h>
#include <string.h>

typedef enum { number, character, string } Type;

typedef struct {
	union {
		int* number;
		char* character;
		char** string;
	} data;
	Type type;
	size_t size;
} Vector;

int sum(Vector* head) {
	size_t i;
	int accum = 0;
	for (i = 0; i < head->size; i++) {
		accum += (head->data.number[i]);
	}

	return accum;
}

int main() {
	// ["A", "B", "C"]
	Vector tmpstrvec0;
	tmpstrvec0.type = string;
	tmpstrvec0.size = 3;
	char* tmpstrvecdata[] = { "A", "B", "C"};
	tmpstrvec0.data.string = tmpstrvecdata;
	// -------------------------------------------------

	// [-2, 5, 43, -44]
	Vector tmpnumvec0;
	tmpnumvec0.type = number;
	tmpnumvec0.size = 4;
	int tmpnumvecdata[] = { -2, 5, 43, -44 };
	tmpnumvec0.data.number = tmpnumvecdata;
	// -------------------------------------------------

	// GEP (["A", "B", "C"], Sum [-2, 5, 43, -44])
	char* First = tmpstrvec0.data.string[sum(&tmpnumvec0)];
	// -------------------------------------------------

	// [['d', 'c', 'a'], "abc"]
	Vector tmpstrvec1;
	tmpstrvec1.type = string;
	tmpstrvec1.size = 2;
	char* tmpstrvecdata1[] = {(char[]){ 'd', 'c', 'a', 0}, "abc"};
	tmpstrvec1.data.string = tmpstrvecdata1;

	// -------------------------------------------------

	// GEP ([['d', 'c', 'a'], "abc"], 0)
	char* Second = tmpstrvec1.data.string[0];
	// -------------------------------------------------

	// Append (
	//	...)
	Vector tmpstrvec2;
	tmpstrvec2.type = character;
	tmpstrvec2.size = 4;
	char* tmpstrvecdata2 = calloc(5, sizeof(char));
	strcpy(tmpstrvecdata2, First);
	strcpy(tmpstrvecdata2+1, Second);

	tmpstrvec2.data.character = tmpstrvecdata2;
	// -------------------------------------------------

	// Println
	puts(tmpstrvec2.data.character);
	// -------------------------------------------------

	return 0;
}
