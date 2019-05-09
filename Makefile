#
#  Copyright 2018 Nalej
# 

include scripts/Makefile.common

.DEFAULT_GOAL := all

# Name of the target applications to be built
APPS=pulsar nalej-pulsar-cli

# Use global Makefile for common targets
export
%:
	@$(MAKE) -f scripts/Makefile.golang $@
