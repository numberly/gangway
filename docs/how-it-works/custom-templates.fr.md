# Modèles Personnalisés

Pour personnaliser les pages HTML rendues par Gangway, vous pouvez fournir un ensemble de modèles personnalisés à utiliser à la place de ceux intégrés.

!!! danger "Important"
    Les données transmises aux modèles peuvent changer entre les versions, et nous ne garantissons pas que nous maintiendrons la compatibilité avec les versions antérieures. Si vous utilisez des modèles personnalisés, une attention particulière doit être portée lors de la mise à niveau de Gangway.**

Pour activer cette fonctionnalité, définissez l'option `customHTMLTemplatesDir` dans le fichier de configuration de Gangway pour un répertoire qui contient les modèles personnalisés suivants :

* home.tmpl : Modèle de la page d'accueil.
* commandline.tmpl : Modèle post-connexion qui liste généralement les commandes nécessaires pour configurer `kubectl`.

Les modèles sont traités en utilisant le [package][0] `html/template` de Go.

[0]: https://golang.org/pkg/html/template/
