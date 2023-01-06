##
## backend
## File description:
## Makefile
##

GO	=	go

NAME1	=	hddManager
NAME2	=	consoManager

SRCDIR1	=	controlHDD
SRCDIR2	=	getConso

SRC1		=	controlHDD.go
SRC2		=	getConso.go

SRC1			:= $(addprefix $(SRCDIR1)/, $(SRC1))
SRC2			:= $(addprefix $(SRCDIR2)/, $(SRC2))

GOFLAGS =	--trimpath --mod=vendor

all: build #lib

build:
	$(GO) mod vendor
	$(GO) build $(GOFLAGS) -o ./$(NAME1) $(SRC1)
	$(GO) build $(GOFLAGS) -o ./$(NAME2) $(SRC2)

# lib:
# 	$(MAKE) -C ./hsmlib

fclean:
	rm -f  $(NAME1) $(NAME2)

re: fclean build #lib