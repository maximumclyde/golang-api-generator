# Golang API Generator

A code generator for rest APIs using gorm and gin.

After cloning, run:

>       make build

to build the binaries. The binary will be created as “api-generator”.

In your project, call

>       ../path/to/api-generator —help

to view the available commands and options.

- **Config**

The bin will look for an optional file named “rest.config.json” which holds the configuration options for the code generation. These options are defined in a json file in the project names “rest.schema.json” defined using the json draft-07 standard. This file can be omitted completely or referenced by a different name using the

>       —config

option and specifying the path to the configuration file.

- **Init**

On a newly starter project, the first thing that should be run is the

>       api-generator init

command. This will create the modules, services, utils and interfaces required to get started.

- **Services**

Services and handlers can be added using the

>       api-generator create service

command. Here you’ll be prompted to input the name of the table in the database. The generator will take care of creating a new migration for the table, a new default model definition, a service for the model and a rest api handler. Along with these, the generator creates or updated en existing store and router. The store saves a reference to the database and the collection of services.

**Use this command only for tables that can be references by a single id field called “id”. This is a simple generator that handles dealing with the boilerplate code that comes with large systems. It’s best for services and handlers to be created manually for more complex tables or tables that use a different convention than the one specified.**

The router takes care of registering routes and exposing them to clients. Upon adding new services these files will be **UPDATED AUTOMATICALLY. Try not to tinker too much in those files.**

The convention is that the router takes all the custom configuration in the beginning and **registers the routes at the end before returning.** In the store, the generator looks up and updates a type called **StoreServices.** This is the collection of all the services for the api, even for those that are not exposed to the api.

To generate services that don’t need to be exposed to external apis, use the command

>       api-generator —no-handler create service

- **Migrations**

A migration file is always created as part of a service, with or without api handlers. In every case, the user will be prompted to input the table name and a default migration will be created using the input value **turned in snake_case.** When created as part of a service, a default migration will write a simple sql statement that creates a table. If needed, the generator can create empty migrations using

>       api-migration create migration

**Careful, this generator is not a migration tool. This is simply a shortcut for writing versioned sql files based on creation timestamp. You’ll need another tool to handle applying and/or rolling back migrations.**
