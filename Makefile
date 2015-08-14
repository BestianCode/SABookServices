TASK = Exporter

all: $(TASK)

%: %.go
#	gccgo -g $< -o $@
	go build -o $@.gl -compiler gc $<

clean:
	rm -f $(TASK) *.gl
