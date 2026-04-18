# ICD Code API

**Fast, developer-friendly ICD-11 search API for real-world applications.**

A simple, fast, and opinionated API for working with ICD-11 diagnostic codes, built to eliminate the friction of using the official WHO interface directly in applications.

Instead of dealing with complex hierarchies, authentication flows, and raw classification data — you get clean, searchable, production-ready responses.


### Why this exists

The official World Health Organization ICD-11 system is powerful, but not designed for real-time application use. 

Building medical software (especially EMRs) introduces common problems:

* **Slow or complex search experience**
* **Hard-to-use hierarchical data**
* **Lack of synonym handling** (“heart attack” vs “myocardial infarction”)
* **Heavy integration overhead** (OAuth, API complexity)
* **Poor UX for autocomplete and clinical workflows**

**ICD Code API** solves this by focusing on developer experience first.


### What makes this different

#### Fast fuzzy search
Search ICD codes using natural language, typos, or partial terms.

* "high sugar" → Type 2 diabetes
* "heart attack" → Acute myocardial infarction
* "bp high" → Essential hypertension

#### Simple JSON response
No complex ontology structures in your frontend.

```json
{
  "code": "5A11",
  "description": "Type 2 diabetes mellitus",
  "version": "ICD-11"
}
```

#### Built for autocomplete UX
Designed for:
* EMR diagnosis selection
* Search-as-you-type interfaces
* Clinical decision support tools

#### No authentication complexity
* No OAuth setup
* No token refresh flows
* Just HTTP requests and results


### Getting Started

For local setup, database initialization, and development workflows, see the **[Getting Started Guide](docs/GETTING_STARTED.md)**.


### Example request
`GET /v1/search?q=heart+attack`

### Example response
```json
[
  {
    "code": "BA41",
    "description": "Acute myocardial infarction",
    "version": "ICD-11"
  }
]
```


### Official WHO Resources

For developers who require official raw datasets or programmatic access directly from the World Health Organization:

* **Official ICD-11 Portal**: [https://icd.who.int/en](https://icd.who.int/en)
* **Official ICD-11 API**: [https://icd.who.int/icdapi](https://icd.who.int/icdapi)
* **Dataset Downloads**: [https://icd.who.int/browse/latest-release/mms/en#/downloaddata](https://icd.who.int/browse/latest-release/mms/en#/downloaddata) *(Requires WHO Registration)*


### Use cases
* Electronic Medical Records (EMR systems)
* Hospital management software
* Clinical decision support tools
* Health tech startups
* Medical research platforms


### Design philosophy
This project is built on a simple idea: **ICD data should be usable, not just accessible.**

So instead of exposing raw classification complexity, we:
* **Flatten** what should be simple
* **Enhance** what should be intuitive
* **Optimize** for real-world workflows, not theoretical correctness


### Roadmap
* [x] Offline dataset support (SQLite)
* [ ] Usage-based ranking improvements
* [ ] Multi-language medical term support
* [ ] Embedding-based semantic search
* [ ] EMR-ready SDKs (Go / Node / Python)


### Contributing

Interested in contributing? Please see our **[Contribution Guidelines](docs/CONTRIBUTING.md)**.


### License
Open-source. Built for the healthcare developer community.

### Why this matters
Healthcare software is often slowed down by tooling that wasn’t designed for real clinicians.

This project aims to make ICD classification: **fast, intuitive, and invisible inside the workflow.**
