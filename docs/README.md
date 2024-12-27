# To-Do List API

Une API robuste pour gérer une to-do list avec des fonctionnalités CRUD pour les tâches. Ce projet diffère du projet `Blog API` par son accent sur l'authentification, la validation avancée des entrées, et la gestion des propriétaires des ressources.

---

## Table des matières
1. Introduction
2. Objectifs du projet
3. Fonctionnalités
4. Technologies utilisées
5. Installation et exécution
6. Structure du projet
7. Documentation API
8. Compétences renforcées
9. Améliorations futures
10. Contributions
11. Licence

---

## Introduction

Cette API RESTful développée en Go permet aux utilisateurs de gérer une liste de tâches. Elle a été conçue pour explorer des techniques avancées comme la validation stricte des entrées, la documentation interactive, et la gestion fine des autorisations.

---

## Objectifs du projet Différences avec le projet `Blog API`

- Apprentissage et application de concepts avancés :
  - Validation stricte des entrées utilisateur.
  - Gestion des relations entre tables, comme les utilisateurs et leurs tâches.
  - Authentification et autorisation sécurisées via cookies de session.

- Fournir une API robuste et maintenable avec une documentation interactive.

---

## Fonctionnalités

- **Gestion des tâches** :
  - Création, consultation, mise à jour, suppression.
  - Filtrage par statut (`complétée`, `en cours`, etc.).

- **Authentification et autorisation** :
  - Middleware `AuthRequired` pour protéger les routes.
  - Middleware `AuthorizeTaskOwnership` pour restreindre l'accès en fonction du propriétaire.

- **Validation stricte des données** :
  - Emails valides, mots de passe sécurisés, et usernames conformes.

- **Documentation interactive** :
  - Documentation complète des endpoints via Swagger.
---

## Technologies utilisées

- **Langage** : Golang
- **Framework HTTP** : Gin
- **ORM** : Gorm
- **Base de données** : SQLite
- **Documentation API** : Swagger

---

## Compétences renforcées

- Gestion des relations entre tables (ex. : utilisateurs et tâches).

- Authentification par cookies de session.

- Validation avancée des entrées utilisateur.

- Gestion fine des autorisations avec des middlewares.

- Implémentation et utilisation de query parameters pour les filtres (ex. : /tasks?status=completed).

- Documentation Swagger interactive.

---

## Documentation API

L'API est entièrement documentée avec Swagger. Pour y accéder :
1. Lancez le serveur.
2. Accédez à l'URL : [http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html).

---

## Améliorations futures

- Ajout d'une pagination pour les tâches.

- Support de bases de données distribuées (PostgreSQL).

- Implémentation de tests d'intégration complets.

- Gestion des sessions avec expiration configurable.

- CI/CD pour les déploiements automatisés.

---

## Contributions

Les contributions sont bienvenues ! Veuillez ouvrir une issue ou une pull request.
