#include <stdint.h>
#include <stdio.h>
#include <stdlib.h>

#define finline __attribute__((always_inline))

typedef struct {
  int32_t length;
  int32_t capacity;
  int32_t width;
  void *data;
} Vector;

finline void writeVectorLength(Vector *v, int32_t length) {
  v->length = length;
}

finline void writeVectorCapacity(Vector *v, int32_t capacity) {
  v->capacity = capacity;
}

finline void writeVectorWidth(Vector *v, int32_t width) { v->width = width; }

finline int32_t readVectorLength(Vector *v) { return v->length; }

finline int32_t readVectorCapacity(Vector *v) { return v->capacity; }

finline int32_t readVectorWidth(Vector *v) { return v->width; }

finline void writeVectorPointer(Vector *v, void *data) { v->data = data; }

finline Vector *createVectorHeader(int32_t length, int32_t capacity,
                                   int32_t width) {
  Vector *v = calloc(1, sizeof(Vector));
  v->length = length;
  v->capacity = capacity;
  v->width = width;
  v->data = NULL;
  return v;
}
