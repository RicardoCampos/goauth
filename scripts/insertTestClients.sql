INSERT INTO public.clients(
	"clientId", "clientSecret", "tokenType", "accessTokenLifetime", "allowedScopes")
	VALUES ('foo_bearer', 'secret', 'Bearer', 0, 'read write');
INSERT INTO public.clients(
	"clientId", "clientSecret", "tokenType", "accessTokenLifetime", "allowedScopes")
	VALUES ('foo_reference', 'secret', 'Reference', 0, 'read write');