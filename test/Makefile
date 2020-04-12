RAND_INPUTS := test.small.dec test.medium.dec test.large.dec
TEXT_INPUTS := test_text.small.dec test_text.medium.dec test_text.large.dec
MIXD_INPUTS := test_mixed.small.dec test_mixed.medium.dec test_mixed.large.dec

ALL_INPUTS := $(RAND_INPUTS) $(TEXT_INPUTS) $(MIXD_INPUTS)
ALL_COMPS  := $(foreach i,$(ALL_INPUTS),$(i).cmp)

LIST := test.list

SML := 1
MED := 100
LRG := 100000

$(LIST): $(ALL_COMPS)
	echo "$(foreach i,$(ALL_COMPS),test/$(i))" > $@

$(ALL_COMPS): %.dec.cmp: %.dec
	lzfse -encode -i $< -o $@

%.small.dec:  COUNT=$(SML)
%.medium.dec: COUNT=$(MED)
%.large.dec:  COUNT=$(LRG)

$(RAND_INPUTS): SOURCE=/dev/urandom
$(TEXT_INPUTS): SOURCE=/usr/share/dict/words

$(RAND_INPUTS) $(TEXT_INPUTS):
	dd if=$(SOURCE) of=$@ bs=1024 count=$(COUNT) 2>/dev/null
$(MIXD_INPUTS):
	dd if=/usr/share/dict/words bs=1024 count=$(COUNT) 2>/dev/null > $@
	dd if=/dev/urandom bs=1024 count=$(COUNT) 2>/dev/null >> $@
	dd if=/usr/share/dict/words bs=1024 count=$(COUNT) 2>/dev/null >> $@
	dd if=/dev/urandom bs=1024 count=$(COUNT) 2>/dev/null >> $@

clean:
	-rm -f $(ALL_COMPS) $(ALL_INPUTS)

.PHONY: clean
