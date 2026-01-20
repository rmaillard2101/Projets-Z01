CC = cc
CFLAGS = -Wall -Wextra -Werror -g

SRC_DIR = src
OBJ_DIR = obj

SRCS = $(SRC_DIR)/main.c $(SRC_DIR)/auth.c $(SRC_DIR)/system.c $(SRC_DIR)/sha256.c
OBJS = $(patsubst $(SRC_DIR)/%.c,$(OBJ_DIR)/%.o,$(SRCS))

TARGET = atm

# Règle par défaut : compile et lie
all: $(TARGET)

# Lien de l'exécutable à partir des fichiers objets
$(TARGET): $(OBJS)
	$(CC) $(CFLAGS) -o $@ $^

# Compilation des fichiers sources en objets
$(OBJ_DIR)/%.o: $(SRC_DIR)/%.c | $(OBJ_DIR)
	$(CC) $(CFLAGS) -c $< -o $@

# Création du dossier obj si inexistant
$(OBJ_DIR):
	mkdir -p $(OBJ_DIR)

# Nettoyer les fichiers objets et l'exécutable
clean:
	rm -rf $(OBJ_DIR) $(TARGET)

.PHONY: all clean
