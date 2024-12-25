# Documentation des Tests

## Introduction
Ce document détaille les tests réalisés pour valider les middlewares `AuthRequired` et `AuthorizeTaskOwnership`, ainsi que les scénarios combinés et avancés pour les routes protégées.

---

## 1. Tests pour AuthRequired

### Cas 1: Accès sans cookie session_token
- **Requête** : Route protégée (`/tasks`, `/tasks/:id`) sans cookie.
- **Attendu** : 
  - Statut : `401 Unauthorized`
  - Réponse : `{"error": "Authentification requise"}`
- **Résultat** : ✅

### Cas 2: Accès avec un session_token invalide
- **Requête** : Ajoute un cookie `session_token=invalid_token`.
- **Attendu** : 
  - Statut : `401 Unauthorized`
  - Réponse : `{"error": "Session invalide"}`
- **Résultat** : ✅

### Cas 3: Accès avec un session_token expiré
- **Préconditions** : Crée une session expirée dans la base de données.
- **Requête** : Ajoute un cookie `session_token` correspondant à cette session.
- **Attendu** : 
  - Statut : `401 Unauthorized`
  - Réponse : `{"error": "Session expirée"}`
- **Résultat** : ✅

### Cas 4: Accès avec un session_token valide
- **Préconditions** : Crée une session valide pour un utilisateur dans la base de données.
- **Requête** : Ajoute un cookie `session_token` valide.
- **Attendu** : 
  - Statut : `200 OK`
  - Réponse : Dépend du contrôleur utilisé après le middleware.
- **Résultat** : ✅

---

## 2. Tests pour AuthorizeTaskOwnership

### Cas 1: Accès sans être authentifié
- **Requête** : Appelle une route protégée (`/tasks/:id`) sans passer par `AuthRequired`.
- **Attendu** : 
  - Statut : `401 Unauthorized`
  - Réponse : `{"error": "Utilisateur non authentifié"}`
- **Résultat** : ✅

### Cas 2: Accès à une tâche inexistante
- **Requête** : Accède à `/tasks/:id` avec un `taskID` inexistant.
- **Attendu** : 
  - Statut : `404 Not Found`
  - Réponse : `{"error": "Tâche non trouvée"}`
- **Résultat** : ✅

### Cas 3: Accès à une tâche qui ne lui appartient pas
- **Préconditions** : 
  - Crée une tâche appartenant à un autre utilisateur.
  - Authentifie l'utilisateur actuel.
- **Requête** : Accède à `/tasks/:id` avec l'ID de la tâche d'un autre utilisateur.
- **Attendu** : 
  - Statut : `401 Unauthorized`
  - Réponse : `{"error": "Action non autorisée"}`
- **Résultat** : ✅

### Cas 4: Accès à une tâche qui lui appartient
- **Préconditions** : Crée une tâche appartenant à l'utilisateur authentifié.
- **Requête** : Accède à `/tasks/:id` avec un ID valide.
- **Attendu** : 
  - Statut : `200 OK`
  - Réponse : La tâche est retournée ou l’opération (mise à jour/suppression) réussit.
- **Résultat** : ✅

---

## 3. Tests des comportements combinés

### Cas 1: Vérifier la chaîne des middlewares
- **Requête** : Accède à une route avec les deux middlewares activés, sans cookie.
- **Attendu** : 
  - Statut : `401 Unauthorized` (dû à `AuthRequired`).
- **Résultat** : ✅

### Cas 2: Session expirée + Tâche inexistante
- **Préconditions** : Utilise un token expiré.
- **Requête** : Accède à `/tasks/:id` avec un `taskID` inexistant.
- **Attendu** : 
  - Statut : `401 Unauthorized` (dû à `AuthRequired`).
- **Résultat** : ✅

### Cas 3: Session valide + Tâche d'un autre utilisateur
- **Préconditions** : 
  - Authentifie l'utilisateur actuel.
  - Crée une tâche pour un autre utilisateur.
- **Requête** : Accède à `/tasks/:id` pour une tâche de l'autre utilisateur.
- **Attendu** : 
  - Statut : `401 Unauthorized` (dû à `AuthorizeTaskOwnership`).
- **Résultat** : ✅

### Cas 4: Session valide + Tâche valide
- **Préconditions** : 
  - Authentifie l'utilisateur actuel.
  - Crée une tâche pour cet utilisateur.
- **Requête** : Accède à `/tasks/:id` pour une tâche valide.
- **Attendu** : 
  - Statut : `200 OK`
  - Réponse : La tâche est retournée ou l’opération réussit.
- **Résultat** : ✅

---

## 4. Scénarios avancés pour les routes protégées

### Cas 1: Créer une tâche avec une session invalide
- **Préconditions** : Utilise un `session_token` invalide.
- **Requête** : `POST /tasks` avec des données valides.
- **Attendu** : 
  - Statut : `401 Unauthorized`
  - Réponse : `{"error": "Authentification requise"}`
- **Résultat** : ✅

### Cas 2: Mettre à jour une tâche d'un autre utilisateur
- **Préconditions** : 
  - Crée une tâche pour un autre utilisateur.
  - Authentifie l'utilisateur actuel.
- **Requête** : `PUT /tasks/:id` avec l'ID de la tâche de l'autre utilisateur.
- **Attendu** : 
  - Statut : `401 Unauthorized`
  - Réponse : `{"error": "Action non autorisée"}`
- **Résultat** : ✅

### Cas 3: Supprimer une tâche sans session active
- **Requête** : `DELETE /tasks/:id` sans cookie.
- **Attendu** : 
  - Statut : `401 Unauthorized`
  - Réponse : `{"error": "Authentification requise"}`
- **Résultat** : ✅

### Cas 4: Accéder à une tâche supprimée
- **Préconditions** : 
  - Crée une tâche.
  - Supprime la tâche.
- **Requête** : Accède à `/tasks/:id` pour une tâche supprimée.
- **Attendu** : 
  - Statut : `404 Not Found`
  - Réponse : `{"error": "Tâche non trouvée"}`
- **Résultat** : ✅

---

## 5. Tests des contextes

### Cas 1: Vérifier `currentUser` injecté dans le contexte
- **Préconditions** : Authentifie un utilisateur valide.
- **Requête** : Appelle une route protégée par `AuthRequired`.
- **Attendu** : Le contrôleur ou le middleware suivant peut accéder à `c.Get("currentUser")` et trouver l'utilisateur.
- **Résultat** : ✅

### Cas 2: Vérifier `task` injectée dans le contexte
- **Préconditions** : 
  - Authentifie un utilisateur.
  - Crée une tâche pour cet utilisateur.
- **Requête** : Appelle une route protégée par `AuthorizeTaskOwnership`.
- **Attendu** : Le contrôleur ou le middleware suivant peut accéder à `c.Get("task")` et trouver la tâche.
- **Résultat** : ✅

---

## Conclusion
Tous les tests ont été réalisés avec succès et les comportements attendus ont été vérifiés. La documentation pourra être mise à jour avec de nouveaux cas au besoin.
