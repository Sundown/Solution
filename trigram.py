dictionary = {}
corpus = ["alzheimers", "Angioedema", "Endometriosis", "Mesothelioma", "Osteoporosis", "Osteosarcoma", "Psoriasis", "Sinusitis", "Trichomonas"]

def trigrams(word):
    for t in [word[i:i+3] for i in range(len(word)-2)]:
        dictionary[t] = dictionary[t] + 1 if t in dictionary else 1

for word in corpus:
    trigrams(word)

for i in sorted(dictionary.items(), key=lambda x: x[1]):
    print(i)
