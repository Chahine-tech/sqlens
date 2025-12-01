🚀 SQLens – Roadmap IA

Objectif : Ajouter une couche d’IA fiable sur SQLens en combinant la puissance du parsing formel + les capacités de raisonnement des LLM.

🎯 Vision

SQLens ne veut pas remplacer ChatGPT.
SQLens veut être un copilote SQL professionnel, fiable, localisable, et intégré aux outils dev/ops.

Grâce au moteur de parsing et d’analyse déjà présent dans le projet, l’IA devient assistée par AST, ce que les LLM généralistes ne peuvent pas offrir.

📌 Roadmap IA (MVP → Avancée)
✅ 1. Intégration LLM (MVP)

 Ajouter un package ai/

 Support OpenAI / Mistral / Claude via HTTP

 Support Ollama (local)

 Fonction générique :
AskLLM(ctx context.Context, prompt string) (string, error)

✅ 2. Explain SQL

 Export AST → JSON

 Prompt IA : “Explique clairement cette requête SQL.”

 Retour structuré : sections, points clés, risques éventuels

 CLI :

sqlens ai explain "SELECT * FROM users"

✅ 3. Correction automatique (“AI Auto-Fix”)

 Récupérer l’erreur exacte du parser/analyzer

 Prompt IA : “Voici la requête et l’erreur, propose une correction.”

 Reparser la suggestion pour validation

 Retourner la meilleure correction valide

 CLI :

sqlens ai fix "SELECT name FROM users u WHERE u.age =="

✅ 4. Optimisation SQL (AI Rewrite)

 Détecter anti-patterns (SELECT *, subqueries inutiles…)

 Prompt IA : “Réécris cette requête pour être plus performante.”

 Comparer AST (diff) entre original et suggestion

 Proposer une version optimisée validée

 CLI :

sqlens ai optimize "SELECT * FROM orders"

✅ 5. Analyse des logs SQL Server

 Parser logs SQL Server déjà supportés

 Identifier requêtes lentes (duration, reads, writes)

 Prompt IA : “Explique pourquoi cette requête est lente.”

 Générer pistes d’optimisation (index, refactor SQL, join hints)

 CLI :

sqlens ai analyze-log slow.log

✅ 6. Caching IA

 Hash du prompt

 Petit cache local .sqlens/cache.db

 Expiration configurable

 Désactivable via variable d’env (SQLENS_AI_CACHE=0)


🧬 8. Fonctionnalités IA avancées (étape suivante)

 Analyse du schéma DB réel pour optimiser mieux

 Index advisor intelligent (LLM + heuristiques)

 Suggestions de partitionnement / clustering

 Mode “review SQL pour PR GitHub”

 Fine-tuning léger local sur AST + requêtes examples