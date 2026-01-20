#include "header.h"
#include <unistd.h>

const char *RECORDS = "./data/records.txt";

int getAccountFromFile(FILE *ptr, char name[50], struct Record *r) {
  return fscanf(ptr, "%d %d %s %d %d/%d/%d %s %d %lf %s", &r->id, &r->userId,
                name, &r->accountNbr, &r->deposit.month, &r->deposit.day,
                &r->deposit.year, r->country, &r->phone, &r->amount,
                r->accountType) != EOF;
}

void saveAccountToFile(FILE *ptr, struct User u, struct Record r) {
  fprintf(ptr, "%d %d %s %d %d/%d/%d %s %d %.2lf %s\n\n", r.id, u.id, u.name,
          r.accountNbr, r.deposit.month, r.deposit.day, r.deposit.year,
          r.country, r.phone, r.amount, r.accountType);
}

void stayOrReturn(int notGood, void f(struct User u), struct User u) {
  int option;
  if (notGood == 0) {
    system("clear");
    printf("\n✖ Record not found!!\n");
  invalid:
    printf("\nEnter 0 to try again, 1 to return to main menu and 2 to exit:");
    scanf("%d", &option);
    if (option == 0)
      f(u);
    else if (option == 1)
      mainMenu(u);
    else if (option == 2)
      exit(0);
    else {
      printf("Insert a valid operation!\n");
      goto invalid;
    }
  } else {
    printf("\nEnter 1 to go to the main menu and 0 to exit:");
    scanf("%d", &option);
  }
  if (option == 1) {
    system("clear");
    mainMenu(u);
  } else {
    system("clear");
    exit(1);
  }
}

void success(struct User u) {
  int option;
  printf("\n✔ Success!\n\n");
invalid:
  printf("Enter 1 to go to the main menu and 0 to exit!\n");
  scanf("%d", &option);
  system("clear");
  if (option == 1) {
    mainMenu(u);
  } else if (option == 0) {
    exit(1);
  } else {
    printf("Insert a valid operation!\n");
    goto invalid;
  }
}

void createNewAcc(struct User u) {
  struct Record r;
  struct Record cr;
  char userName[50];
  FILE *pf = fopen(RECORDS, "a+");
  if (!pf) {
    printf("Erreur ouverture fichier\n");
    return;
  }

noAccount:
  system("clear");
  printf("\t\t\t===== New record =====\n");

  printf("\nEnter today's date(mm/dd/yyyy):");
  scanf("%d/%d/%d", &r.deposit.month, &r.deposit.day, &r.deposit.year);
  printf("\nEnter the account number:");
  scanf("%d", &r.accountNbr);

  rewind(pf);

  while (getAccountFromFile(pf, userName, &cr)) {
    if (cr.userId == u.id && cr.accountNbr == r.accountNbr) {
      printf("✖ This Account already exists for this user\n\n");
      goto noAccount;
    }
  }

  printf("\nEnter the country:");
  scanf("%s", r.country);
  printf("\nEnter the phone number:");
  scanf("%d", &r.phone);
  printf("\nEnter amount to deposit: $");
  scanf("%lf", &r.amount);
  printf("\nChoose the type of account:\n\t-> saving\n\t-> current\n\t-> "
         "fixed01(for 1 year)\n\t-> fixed02(for 2 years)\n\t-> fixed03(for 3 "
         "years)\n\n\tEnter your choice:");
  scanf("%s", r.accountType);

  if (strcmp(r.accountType, "current") != 0 &&
      strcmp(r.accountType, "saving") != 0 &&
      strcmp(r.accountType, "fixed01") != 0 &&
      strcmp(r.accountType, "fixed02") != 0 &&
      strcmp(r.accountType, "fixed03") != 0) {

    printf("✖ Invalid account type!\n");
    fclose(
        pf); // N'oublie pas de fermer le fichier avant d'appeler stayOrReturn
    stayOrReturn(0, createNewAcc, u);
    return; // Important pour ne pas continuer la fonction après
  }
  int lastId = 0;
  rewind(pf);
  while (getAccountFromFile(pf, userName, &cr)) {
    if (cr.id > lastId)
      lastId = cr.id;
  }
  r.id = lastId + 1;

  // Affectation correcte de userId
  r.userId = u.id;

  fseek(pf, 0, SEEK_END);

  saveAccountToFile(pf, u, r);

  fclose(pf);
  success(u);
}

void updateAccountInfo(struct User u) {
  FILE *fp = fopen(RECORDS, "r");
  FILE *temp = fopen("./data/temp.txt", "w");
  struct Record r;
  char userName[50];
  int found = 0;
  int accNum, choice;

  if (!fp || !temp) {
    printf("Error opening file.\n");
    exit(1);
  }

  system("clear");
  printf("\n\t\t=== Update Account Information ===\n");

  printf("Enter the account NUMBER you want to update: ");
  scanf("%d", &accNum);

  while (getAccountFromFile(fp, userName, &r)) {
    if (strcmp(userName, u.name) == 0 && r.accountNbr == accNum) {
      found = 1;

      printf("\nWhich field do you want to update?\n");
      printf("1. Country\n");
      printf("2. Phone number\n");
      printf("Enter your choice: ");
      scanf("%d", &choice);

      if (choice == 1) {
        printf("Enter new country: ");
        scanf("%s", r.country);
      } else if (choice == 2) {
        printf("Enter new phone number: ");
        scanf("%d", &r.phone);
      } else {
        printf("Invalid choice.\n");
        fclose(fp);
        fclose(temp);
        remove("./data/temp.txt");
        return;
      }
    }

    fprintf(temp, "%d %d %s %d %d/%d/%d %s %d %.2lf %s\n\n", r.id, r.userId,
            userName, r.accountNbr, r.deposit.month, r.deposit.day,
            r.deposit.year, r.country, r.phone, r.amount, r.accountType);
  }

  fclose(fp);
  fclose(temp);

  remove(RECORDS);
  rename("./data/temp.txt", RECORDS);

  if (found) {
    printf("\n✔ Account information updated successfully.\n");
    success(u);
  } else {
    printf("\n✖ Account with number %d not found or does not belong to you.\n",
           accNum);
    stayOrReturn(0, updateAccountInfo, u);
  }
}

void checkOneAccount(struct User u) {
  int accountId;
  char ownerName[100];
  struct Record r;
  FILE *fp = fopen(RECORDS, "r");

  if (fp == NULL) {
    perror("Error opening records file");
    exit(1);
  }

  system("clear");
  printf("\n\t\t====== View Specific Account ======\n\n");
  printf("Enter the account number to view: ");
  scanf("%d", &accountId);

  int found = 0;

  while (getAccountFromFile(fp, ownerName, &r)) {
    if (strcmp(ownerName, u.name) == 0 && r.accountNbr == accountId) {
      found = 1;
      printf("\n\t✓ Account found!\n\n");
      printf("Account number: %d\n", r.accountNbr);
      printf("Deposit Date : %d/%d/%d\n", r.deposit.day, r.deposit.month,
             r.deposit.year);
      printf("Country      : %s\n", r.country);
      printf("Phone        : %d\n", r.phone);
      printf("Amount       : $%.2lf\n", r.amount);
      printf("Type         : %s\n", r.accountType);

      double interest = 0.0;

      if (strcmp(r.accountType, "saving") == 0) {
        interest = (r.amount * 0.07) / 12;
        printf("You will get $%.2lf as interest on day %d of every month.\n",
               interest, r.deposit.day);
      } else if (strcmp(r.accountType, "fixed01") == 0) {
        interest = (r.amount * 0.04);
        printf("You will get $%.2lf as interest on day %d/%d of %d.\n",
               interest, r.deposit.day, r.deposit.month, (r.deposit.year + 1));
      } else if (strcmp(r.accountType, "fixed02") == 0) {
        interest = (r.amount * 0.05) * 2;
        printf("You will get $%.2lf as interest on day %d/%d of %d.\n",
               interest, r.deposit.day, r.deposit.month, (r.deposit.year + 2));
      } else if (strcmp(r.accountType, "fixed03") == 0) {
        interest = (r.amount * 0.08) * 3;
        printf("You will get $%.2lf as interest on day %d/%d of %d.\n",
               interest, r.deposit.day, r.deposit.month, (r.deposit.year + 3));
      } else if (strcmp(r.accountType, "current") == 0) {
        printf("You will not get interests because the account is of type "
               "current.\n");
      } else {
        printf("❗ Unknown account type: %s\n", r.accountType);
      }

      break;
    }
  }

  fclose(fp);

  if (!found) {
    printf("\n✖ Account not found for this user.\n");
  }

  success(u);
}

void checkAllAccounts(struct User u) {
  char userName[100];
  struct Record r;

  FILE *pf = fopen(RECORDS, "r");

  system("clear");
  printf("\t\t====== All accounts from user, %s =====\n\n", u.name);
  while (getAccountFromFile(pf, userName, &r)) {
    if (strcmp(userName, u.name) == 0) {
      printf("_____________________\n");
      printf("\nAccount number:%d\nDeposit Date:%d/%d/%d \ncountry:%s \nPhone "
             "number:%d \nAmount deposited: $%.2f \nType Of Account:%s\n",
             r.accountNbr, r.deposit.day, r.deposit.month, r.deposit.year,
             r.country, r.phone, r.amount, r.accountType);
    }
  }
  fclose(pf);
  success(u);
}

void makeTransaction(struct User u) {
  FILE *fp = fopen(RECORDS, "r");
  FILE *temp = fopen("./data/temp.txt", "w");
  struct Record r;
  int senderAccNum, receiverAccNum;
  double amount;
  int senderFound = 0, receiverFound = 0;
  struct Record sender, receiver;
  char name[100];
  char receiverName[100];
  int choice;

  if (!fp || !temp) {
    perror("Erreur lors de l’ouverture des fichiers.");
    if (fp)
      fclose(fp);
    if (temp)
      fclose(temp);
    exit(1);
  }

  system("clear");
  printf("\n\t\t===== Opération bancaire =====\n\n");
  printf("Choisissez une opération:\n");
  printf("1 - Transfert\n");
  printf("2 - Dépôt\n");
  printf("3 - Retrait\n");
  printf("Votre choix : ");
  scanf("%d", &choice);

  printf("Entrez le numéro de votre compte : ");
  scanf("%d", &senderAccNum);

  // Recherche du compte utilisateur
  while (getAccountFromFile(fp, name, &r)) {
    if (strcmp(name, u.name) == 0 && r.accountNbr == senderAccNum) {
      sender = r;
      senderFound = 1;
      strcpy(sender.name, name);
      break;
    }
  }

  if (!senderFound) {
    fclose(fp);
    fclose(temp);
    printf("\n✖ Compte introuvable ou non associé à votre utilisateur.\n");
    stayOrReturn(0, makeTransaction, u);
    return;
  }

  if (choice == 1) {
    // Transfert
    printf("Entrez le nom du destinataire : ");
    scanf("%s", receiverName);
    printf("Entrez le numéro du compte destinataire : ");
    scanf("%d", &receiverAccNum);

    rewind(fp);
    receiverFound = 0;
    while (getAccountFromFile(fp, name, &r)) {
      if (strcmp(name, receiverName) == 0 && r.accountNbr == receiverAccNum) {
        receiver = r;
        receiverFound = 1;
        strcpy(receiver.name, name);
        break;
      }
    }

    if (!receiverFound) {
      fclose(fp);
      fclose(temp);
      printf("\n✖ Aucun compte #%d appartenant à %s trouvé.\n", receiverAccNum,
             receiverName);
      stayOrReturn(0, makeTransaction, u);
      return;
    }

    if (sender.accountNbr == receiver.accountNbr &&
        strcmp(sender.name, receiver.name) == 0) {
      fclose(fp);
      fclose(temp);
      printf("\n✖ Vous ne pouvez pas transférer de l'argent vers votre propre "
             "compte.\n");
      stayOrReturn(0, makeTransaction, u);
      return;
    }

    printf("Entrez le montant à transférer : $");
    scanf("%lf", &amount);

    if (sender.amount < amount) {
      fclose(fp);
      fclose(temp);
      printf("\n✖ Fonds insuffisants. Solde actuel : $%.2lf\n", sender.amount);
      stayOrReturn(0, makeTransaction, u);
      return;
    }

    if (strcmp(sender.accountType, "fixed01") == 0 ||
        strcmp(sender.accountType, "fixed02") == 0 ||
        strcmp(sender.accountType, "fixed03") == 0) {
      fclose(fp);
      fclose(temp);
      printf("\n✖ Ce compte est à terme (fixed), les transferts sont "
             "interdits.\n");
      stayOrReturn(0, makeTransaction, u);
      return;
    }

    // Effectuer transfert
    sender.amount -= amount;
    receiver.amount += amount;

  } else if (choice == 2) {
    // Dépôt
    printf("Entrez le montant à déposer : $");
    scanf("%lf", &amount);

    if (amount <= 0) {
      fclose(fp);
      fclose(temp);
      printf("\n✖ Le montant de dépôt doit être positif.\n");
      stayOrReturn(0, makeTransaction, u);
      return;
    }

    if (strcmp(sender.accountType, "fixed01") == 0 ||
        strcmp(sender.accountType, "fixed02") == 0 ||
        strcmp(sender.accountType, "fixed03") == 0) {
      fclose(fp);
      fclose(temp);
      printf("\n✖ Ce compte est à terme (fixed), les dépôts sont interdits.\n");
      stayOrReturn(0, makeTransaction, u);
      return;
    }

    sender.amount += amount;

  } else if (choice == 3) {
    // Retrait
    printf("Entrez le montant à retirer : $");
    scanf("%lf", &amount);

    if (amount <= 0) {
      fclose(fp);
      fclose(temp);
      printf("\n✖ Le montant de retrait doit être positif.\n");
      stayOrReturn(0, makeTransaction, u);
      return;
    }

    if (strcmp(sender.accountType, "fixed01") == 0 ||
        strcmp(sender.accountType, "fixed02") == 0 ||
        strcmp(sender.accountType, "fixed03") == 0) {
      fclose(fp);
      fclose(temp);
      printf(
          "\n✖ Ce compte est à terme (fixed), les retraits sont interdits.\n");
      stayOrReturn(0, makeTransaction, u);
      return;
    }

    if (sender.amount < amount) {
      fclose(fp);
      fclose(temp);
      printf("\n✖ Fonds insuffisants. Solde actuel : $%.2lf\n", sender.amount);
      stayOrReturn(0, makeTransaction, u);
      return;
    }

    sender.amount -= amount;

  } else {
    fclose(fp);
    fclose(temp);
    printf("\n✖ Option invalide.\n");
    stayOrReturn(0, makeTransaction, u);
    return;
  }

  rewind(fp);

  // Réécriture des comptes dans le fichier temporaire
  while (getAccountFromFile(fp, name, &r)) {
    struct User currentUser;
    currentUser.id = r.userId;
    strcpy(currentUser.name, name);
    strcpy(currentUser.password, "N/A");

    if (r.accountNbr == sender.accountNbr && strcmp(name, sender.name) == 0) {
      saveAccountToFile(temp, currentUser, sender);
    } else if (choice == 1 && r.accountNbr == receiver.accountNbr &&
               strcmp(name, receiver.name) == 0) {
      saveAccountToFile(temp, currentUser, receiver);
    } else {
      saveAccountToFile(temp, currentUser, r);
    }
  }

  fclose(fp);
  fclose(temp);

  remove(RECORDS);
  rename("./data/temp.txt", RECORDS);

  if (choice == 1) {
    printf("\n✔ Transfert de $%.2lf effectué avec succès !\n", amount);
    printf("À : %s (compte #%d)\n", receiver.name, receiver.accountNbr);
  } else if (choice == 2) {
    printf("\n✔ Dépôt de $%.2lf effectué avec succès !\n", amount);
  } else if (choice == 3) {
    printf("\n✔ Retrait de $%.2lf effectué avec succès !\n", amount);
  }

  printf("Nouveau solde : $%.2lf\n", sender.amount);
  success(u);
}

void deleteAccount(struct User u) {
  FILE *fp = fopen(RECORDS, "r");
  FILE *temp = fopen("./data/temp.txt", "w");
  struct Record r;
  char name[50];
  int targetAccNum;
  int found = 0;

  if (!fp || !temp) {
    perror("✖ Erreur lors de l’ouverture des fichiers.");
    if (fp)
      fclose(fp);
    if (temp)
      fclose(temp);
    return;
  }

  system("clear");
  printf("\n\t\t===== ✘ Suppression de compte ✘ =====\n\n");
  printf("Utilisateur connecté : %s\n", u.name);
  printf("--------------------------------------------------\n");

  printf("Entrez le numéro du compte à supprimer : ");
  scanf("%d", &targetAccNum);

  // Parcours et écriture dans le fichier temporaire (en excluant le compte
  // ciblé)
  while (getAccountFromFile(fp, name, &r)) {
    if (r.accountNbr == targetAccNum && strcmp(name, u.name) == 0) {
      found = 1; // Le compte est bien à supprimer
      continue;  // Ne pas le recopier
    }

    // Copier les autres comptes
    struct User dummyUser;
    dummyUser.id = r.userId;
    strcpy(dummyUser.name, name);
    strcpy(dummyUser.password, "N/A");
    saveAccountToFile(temp, dummyUser, r);
  }

  fclose(fp);
  fclose(temp);

  if (found) {
    remove(RECORDS);
    rename("./data/temp.txt", RECORDS);
    printf("\n✔ Compte #%d supprimé avec succès.\n", targetAccNum);
  } else {
    remove("./data/temp.txt");
    printf("\n✖ Aucun compte #%d associé à l’utilisateur %s trouvé.\n",
           targetAccNum, u.name);
  }

  success(u); // Retour au menu principal
}

void transferAccountOwnership(struct User u) {
  FILE *fp = fopen(RECORDS, "r");
  FILE *temp = fopen("./data/temp.txt", "w");
  struct Record r;
  char name[100];
  int targetAccNum;
  char newOwnerName[50];
  int found = 0, newOwnerFound = 0;
  struct User newOwner;

  if (!fp || !temp) {
    perror("Erreur d'ouverture de fichier");
    if (fp)
      fclose(fp);
    if (temp)
      fclose(temp);
    return;
  }

  system("clear");
  printf("\n\t\t===== Transfert de propriété d'un compte =====\n\n");

  printf("Entrez le numéro du compte à transférer : ");
  scanf("%d", &targetAccNum);

  printf("Entrez le nom du nouveau propriétaire : ");
  scanf("%s", newOwnerName);

  // Recherche du nouveau propriétaire dans le fichier users
  FILE *users = fopen("./data/users.txt", "r");
  if (!users) {
    perror("Erreur lors de l’ouverture de la base utilisateurs");
    fclose(fp);
    fclose(temp);
    return;
  }

  char uname[50], pass[50];
  int uid;
  while (fscanf(users, "%d %s %s\n", &uid, uname, pass) == 3) {
    if (strcmp(uname, newOwnerName) == 0) {
      newOwnerFound = 1;
      newOwner.id = uid;
      strcpy(newOwner.name, uname);
      strcpy(newOwner.password, pass); // inutile ici mais OK
      break;
    }
  }
  fclose(users);

  if (!newOwnerFound) {
    fclose(fp);
    fclose(temp);
    printf("\n✖ L'utilisateur '%s' n'existe pas.\n", newOwnerName);
    stayOrReturn(0, transferAccountOwnership, u);
    return;
  }

  if (strcmp(u.name, newOwner.name) == 0) {
    fclose(fp);
    fclose(temp);
    printf("\n✖ Vous ne pouvez pas transférer un compte à vous-même.\n");
    stayOrReturn(0, transferAccountOwnership, u);
    return;
  }

  // 🔧 MODIFICATION : on prépare à vérifier les conflits
  int newOwnerAccountNumbers[1000];
  int newOwnerAccountCount = 0;

  FILE *fpScan = fopen(RECORDS, "r");
  struct Record tempR;
  char tempName[100];
  while (getAccountFromFile(fpScan, tempName, &tempR)) {
    if (strcmp(tempName, newOwner.name) == 0) {
      newOwnerAccountNumbers[newOwnerAccountCount++] = tempR.accountNbr;
    }
  }
  fclose(fpScan);

  // 🔁 Parcours des comptes dans RECORDS
  while (getAccountFromFile(fp, name, &r)) {
    if (r.accountNbr == targetAccNum && strcmp(name, u.name) == 0) {
      found = 1;

      r.userId = newOwner.id;
      strcpy(name, newOwner.name);

      // 🔧 MODIFICATION : vérifie conflit de numéro
      int conflict = 0;
      for (int i = 0; i < newOwnerAccountCount; i++) {
        if (newOwnerAccountNumbers[i] == targetAccNum) {
          conflict = 1;
          break;
        }
      }

      if (conflict) {
        int oldNbr = r.accountNbr;
        int newNbr = 1;
        while (1) {
          int alreadyUsed = 0;
          for (int i = 0; i < newOwnerAccountCount; i++) {
            if (newOwnerAccountNumbers[i] == newNbr) {
              alreadyUsed = 1;
              break;
            }
          }
          if (!alreadyUsed)
            break;
          newNbr++;
        }
        r.accountNbr = newNbr;
        printf("\nℹ️ Conflit détecté : %s possède déjà un compte #%d. Le compte "
               "transféré sera renommé en #%d.\n",
               newOwner.name, oldNbr, r.accountNbr);
      }

      saveAccountToFile(temp, newOwner, r);
    } else {
      struct User dummy;
      dummy.id = r.userId;
      strcpy(dummy.name, name);
      strcpy(dummy.password, "N/A");
      saveAccountToFile(temp, dummy, r);
    }
  }

  fclose(fp);
  fclose(temp);

  if (found) {
    remove(RECORDS);
    rename("./data/temp.txt", RECORDS);
    printf("\n✔ Compte #%d transféré avec succès à %s.\n", targetAccNum,
           newOwner.name);
  } else {
    remove("./data/temp.txt");
    printf("\n✖ Aucun compte #%d appartenant à %s trouvé.\n", targetAccNum,
           u.name);
  }

  success(u);
}
