#include <stdint.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

struct Date {
  int month, day, year;
};

// all fields for each record of an account
struct Record {
  int id;
  int userId;
  char name[100];
  char country[100];
  int phone;
  char accountType[10];
  int accountNbr;
  double amount;
  struct Date deposit;
  struct Date withdraw;
};

struct User {
  int id;
  char name[50];
  char password[65];
};

// authentication functions
void loginMenu(char a[50], char pass[65]);
void registerMenu(char a[50], char pass[65]);
void Registration(struct User *u);
void getPassword(struct User u, char *buff, size_t buff_size);
void getId(struct User *u);

// system function
void mainMenu(struct User u);
void createNewAcc(struct User u);
void updateAccountInfo(struct User u);
void checkOneAccount(struct User u);
void checkAllAccounts(struct User u);
void makeTransaction(struct User u);
void deleteAccount(struct User u);
void transferAccountOwnership(struct User u);