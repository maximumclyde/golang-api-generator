# Golang API Generator

A code generator for rest APIs using gorm and gin.

After cloning, run:

```bash
make build
```

to build the binaries. The binary will be created as “api-generator”.

In your project, call

```bash
../path/to/api-generator --help
```

to view the available commands and options.

- **Config**

The bin will look for an **optional** file named “generator.config.json” which holds the configuration options for the code generation. These options are defined in a json file in the project names “rest.schema.json” defined using the json draft-07 standard. This file can be omitted completely or referenced by a different name using the

```bash
--config
```

option and specifying the path to the configuration file.

To create a default configuration file, run

```bash
api-generator [--config="./some/path"] create config
```

- **Init**

On a newly starter project, the first thing that should be run is

```bash
api-generator init
```

This will create the modules, services, utils and interfaces required to get started.

- **Services**

Services and handlers can be added using the

```bash
api-generator create service
```

command. Here you’ll be prompted to input the name of the table in the database. The generator will take care of creating:

- a new migration for the table
- a new default model definition
- a service for the model
- a seeder for the model
- a rest api handler (if the --no-handler option is false)

Along with these, the generator creates or updated the existing store, router and seeder main file. The models have default decoration for json conversion, gorm usage and **[go-faker](https://github.com/go-faker/faker)** for seeders

If a service needs to be created for a table that is not referenced by a single primary key, so it doesn't follow the standard convention, run:

```bash
api-generator --custom create service
```

This command will still create all the mentioned files, however, it's up to you to manually implement all the functions and business logic.

The router takes care of registering routes and exposing them to clients. Upon adding new services these files will be **UPDATED AUTOMATICALLY. Try not to tinker too much in those files.**

The convention is that the router takes all the custom configuration in the beginning and **registers the routes at the end before returning.** In the store, the generator looks up and updates a type called **StoreServices.** This is the collection of all the services in the application, even for those that are not exposed to the api. The seeder file gets updated automatically. New seeders are **appended at the end of the existing file**

To generate services that don’t need to be exposed to external apis, use the command

```bash
api-generator --no-handler create service
```

- **Migrations**

A migration file is always created as part of a service, with or without api handlers. In every case, the user will be prompted to input the table name and a default migration will be created using the input value **turned in snake_case.** When created as part of a service, a default migration will write a simple sql statement that creates a table. If needed, the generator can create empty migrations using

```bash
api-migration create migration
```

**Careful, this generator is not a migration tool. This is simply a shortcut for writing versioned sql files based on creation timestamp. You’ll need another tool to handle applying and/or rolling back migrations.**
