# stringservice

Basic go-kit service, based on go-kit stringservice examples, meant to be deployed on kubernetes.
Does lowercase and count, depends on lowercase service for lowercase endpoint.
Leverages kubernetes service discovery to discover lowercase service (http://lowercase:80)
