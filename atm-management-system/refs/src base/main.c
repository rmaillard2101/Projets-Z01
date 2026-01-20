#include "header.h"

void mainMenu(struct User u) {
  int option;
  system("clear");
  printf("\n\n\t\t======= ATM =======\n\n");
  printf("\n\t\t-->> Feel free to choose one of the options below <<--\n");
  printf("\n\t\t[1]- Create a new account\n");
  printf("\n\t\t[2]- Update account information\n");
  printf("\n\t\t[3]- Check accounts\n");
  printf("\n\t\t[4]- Check list of owned account\n");
  printf("\n\t\t[5]- Make Transaction\n");
  printf("\n\t\t[6]- Remove existing account\n");
  printf("\n\t\t[7]- Transfer ownership\n");
  printf("\n\t\t[8]- Exit\n");
  // printf("%d\n", u.id);
  // printf("%s\n", u.name);
  // printf("DEBUG: u.id = %d, u.name = %s, u.password = %s\n", u->id,
  // u->name,u->password);
  scanf("%d", &option);

  switch (option) {
  case 1:
    createNewAcc(u);
    break;
  case 2:
    updateAccountInfo(u);
    break;
  case 3:
    checkOneAccount(u);
    break;
  case 4:
    checkAllAccounts(u);
    break;
  case 5:
    makeTransaction(u);
    break;
  case 6:
    deleteAccount(u);
    break;
  case 7:
    transferAccountOwnership(u);
    break;
  case 8:
    exit(1);
    break;
  default:
    printf("Invalid operation!\n");
  }
};

void initMenu(struct User *u) {
  int r = 0;
  int option;
  system("clear");
  printf("\n\n\t\t======= ATM =======\n");
  printf("\n\t\t-->> Feel free to login / register :\n");
  printf("\n\t\t[1]- login\n");
  printf("\n\t\t[2]- register\n");
  printf("\n\t\t[3]- exit\n");
  while (!r) {
    scanf("%d", &option);
    switch (option) {
    case 1:
      char pass[50];
      loginMenu(u->name, u->password);
      // printf("User entered: %s / %s\n", u->name, u->password);
      getPassword(*u, pass, sizeof(pass));
      // printf("Password found: %s \n", pass);
      if (strcmp(u->password, pass) == 0) {
        printf("\n\nPassword Match!");
        getId(u);
      } else {
        printf("\nWrong password!! or User Name\n");
        exit(1);
      }
      r = 1;
      break;
    case 2:
      Registration(u);
      r = 1;
      break;
    case 3:
      exit(1);
      break;
    default:
      printf("Insert a valid operation!\n");
    }
  }
};

int main() {
  struct User u;

  initMenu(&u);
  mainMenu(u);
  return 0;
}
