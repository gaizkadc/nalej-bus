#
#  Copyright 2018 Nalej
# 

# Name of the target applications to be built
APPS=pulsar nalej-pulsar-cli

# Use global Makefile for common targets
export
%:
	$(MAKE) -f Makefile.golang $@
