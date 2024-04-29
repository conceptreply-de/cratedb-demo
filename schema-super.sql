-- Creating analyzer for vin search requires superuser privileges
CREATE ANALYZER three_gram_analyzer (
	TOKENIZER three_gram with ( type = 'ngram', min_gram  = 3, max_gram = 3)
);

