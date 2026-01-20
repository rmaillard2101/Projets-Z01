#include "header.h"
#include "sha256.h"
#include <termios.h>

char *USERS = "./data/users.txt";

void loginMenu(char a[50], char pass[65]) {
  struct termios oflags, nflags;

  system("clear");
  printf("\n\n\n\t\t\t\t   Bank Management System\n\t\t\t\t\t User Login:");
  scanf("%s", a);

  // disabling echo
  tcgetattr(fileno(stdin), &oflags);
  nflags = oflags;
  nflags.c_lflag &= ~ECHO;
  nflags.c_lflag |= ECHONL;

  if (tcsetattr(fileno(stdin), TCSANOW, &nflags) != 0) {
    perror("tcsetattr");
    return exit(1);
  }
  printf("\n\n\n\n\n\t\t\t\tEnter the password to login:");
  scanf("%s", pass);

  // restore terminal
  if (tcsetattr(fileno(stdin), TCSANOW, &oflags) != 0) {
    perror("tcsetattr");
    return exit(1);
  }
};

void Registration(struct User *u) {
  FILE *fp;
  struct User temp;
  int lastId = -1;
  char confirmPassword[65];
  struct termios oflags, nflags;

  system("clear");
  printf("\n\n\t\t\t\t   User Registration");

  printf("\n\nEnter a username: ");
  scanf("%s", u->name);

  // Disable terminal echo
  tcgetattr(fileno(stdin), &oflags);
  nflags = oflags;
  nflags.c_lflag &= ~ECHO;
  nflags.c_lflag |= ECHONL;
  if (tcsetattr(fileno(stdin), TCSANOW, &nflags) != 0) {
    perror("tcsetattr");
    exit(1);
  }

  printf("Enter a password: ");
  scanf("%s", u->password);

  printf("Confirm your password: ");
  scanf("%s", confirmPassword);

  // Re-enable terminal echo
  if (tcsetattr(fileno(stdin), TCSANOW, &oflags) != 0) {
    perror("tcsetattr");
    exit(1);
  }

  // Check if passwords match
  if (strcmp(u->password, confirmPassword) != 0) {
    printf("\nPasswords do not match. Registration aborted.\n");
    exit(1);
  }

  // Check for existing username and find last used ID
  if ((fp = fopen(USERS, "r")) != NULL) {
    while (fscanf(fp, "%d %s %s", &temp.id, temp.name, temp.password) != EOF) {
      if (strcmp(temp.name, u->name) == 0) {
        printf("\nUsername already exists.\n");
        fclose(fp);
        exit(1);
      }
      if (temp.id > lastId)
        lastId = temp.id;
    }
    fclose(fp);
  }

  u->id = lastId + 1;

  fp = fopen(USERS, "a");
  if (fp == NULL) {
    perror("Error opening file");
    exit(1);
  }

  char hashed[65];
  hash_password(u->password, hashed);
  fprintf(fp, "%d %s %s\n", u->id, u->name, hashed);
  fclose(fp);

  printf("\nRegistration successful! Welcome, %s!\n", u->name);
}

void getPassword(struct User u, char *buff, size_t buff_size) {
  FILE *fp;
  struct User userChecker;

  if ((fp = fopen("./data/users.txt", "r")) == NULL) {
    printf("Error! opening file");
    exit(1);
  }

  while (fscanf(fp, "%d %s %s", &userChecker.id, userChecker.name,
                userChecker.password) != EOF) {
    if (strcmp(userChecker.name, u.name) == 0) {
      strncpy(buff, userChecker.password, buff_size - 1);
      buff[buff_size - 1] = '\0'; // pour sécurité
      fclose(fp);
      return;
    }
  }

  buff[0] = '\0';
  fclose(fp);
}

void getId(struct User *u) {
  FILE *fp;
  struct User userChecker;

  if ((fp = fopen("./data/users.txt", "r")) == NULL) {
    printf("Error! opening file");
    exit(1);
  }

  while (fscanf(fp, "%d %s %s", &userChecker.id, userChecker.name,
                userChecker.password) != EOF) {
    if (strcmp(userChecker.name, u->name) == 0) {
      u->id = userChecker.id;
      fclose(fp);
      return;
    }
  }

  fclose(fp);
  printf("Erreur : ID utilisateur introuvable.\n");
  exit(1);
}

// Convertir le hash binaire en chaîne hexadécimale
void sha256_to_hex(const uint8_t *hash, char *output) {
  for (int i = 0; i < 32; i++) {
    sprintf(output + i * 2, "%02x", hash[i]);
  }
  output[64] = '\0';
}

// Fonction pratique pour résumer le hash en une seule étape
void hash_password(const char *password, char output[65]) {
  SHA256_CTX ctx;
  uint8_t hash[32];
  sha256_init(&ctx);
  sha256_update(&ctx, (const uint8_t *)password, strlen(password));
  sha256_final(&ctx, hash);
  sha256_to_hex(hash, output);
}
