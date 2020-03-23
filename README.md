Classpath Builder
==================
This tool helps to create a classpath from differen project folders for
complex java projects.

Usage
-----
It is possible to configure the project directory, the path of  
subproject jar files  and a list of excluded project directories.
Directories with a dot as the first character in the name are not included.

```shell script
Usage of cpbuilder:
  -dir string
        Directory with projects / cartridges
  -excludes string
        String with excluded directories
  -path string
        Sub path with jar file
```

It is also possible to configure the exclude list or the library path
over environment variables.

| Name           | Description                           | Example          |
|----------------|---------------------------------------|------------------|
| `LIBPATH`      | Path with jar files                   | build/libs       |
| `PATHEXCLUDES` | List with exluded project directories | build,config,bin |

Example
-------
Project structure
```
/user/name/project
   |- .gradle
   |- bin
   |- build
   |- subprj1
   |- subprj2
   |- gradle
   |- build.gradle
   ...
```

Skript
```shell script
export PATHEXCLUDES="bin,build,gradle"
CLASSPATH=`cpbuilder -dir /user/name/project -path target/libs

echo $CLASSPATH
```

Result
```
/user/name/project/subprj1/target/libs/*:/user/name/project/subprj2/target/libs/*
```

License
------------
Copyright 2014-2020 Matthias Raab.

Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with the License. You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.

