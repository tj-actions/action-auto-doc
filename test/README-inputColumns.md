# test

## Inputs

<!-- AUTO-DOC-INPUT:START - Do not remove or modify this section -->

|     INPUT      |  TYPE  |                                     DESCRIPTION                                     |
|----------------|--------|-------------------------------------------------------------------------------------|
|    base_sha    | string |               Specify a base commit SHA<br>used for comparing changes               |
|     files      | string | Check for changes using only<br>this list of files (Defaults<br>to the entire repo) |
|      path      | string |   Specify a relative path under<br>$GITHUB_WORKSPACE to locate the repository<br>   |
| path_separator | string |                                   Path separator                                    |
|   separator    | string |                        Split character for array output<br>                         |
|      sha       | string |             Specify a current commit SHA<br>used for comparing changes              |
|     token      | string |                Github token or Repo Scoped<br>Personal Access Token                 |

<!-- AUTO-DOC-INPUT:END -->

## Outputs

<!-- AUTO-DOC-OUTPUT:START - Do not remove or modify this section -->

|             OUTPUT             |  TYPE  |                                       DESCRIPTION                                       |
|--------------------------------|--------|-----------------------------------------------------------------------------------------|
|          added_files           | string |                                  List of added files.                                   |
| all_changed_and_modified_files | string |                             List of all changed files.<br>                              |
|       all_modified_files       | string |                     List of all copied modified<br>and added files.                     |
|          any_changed           | string |   Return true only when any<br>files provided using the files<br>input have changed.    |
|          any_deleted           | string | Return true only when any<br>files provided using the files<br>input have been deleted. |
|          copied_files          | string |                                  List of copied files.                                  |
|         deleted_files          | string |                                 List of deleted files.                                  |
|         modified_files         | string |                                 List of modified files.                                 |
|          only_changed          | string |      Return true when all files<br>provided using the files input<br>have changed.      |
|          only_deleted          | string |   Return true when all files<br>provided using the files input<br>have been deleted.    |
|      other_deleted_files       | string |            Return list of deleted files<br>not listed in the files<br>input.            |
|         renamed_files          | string |                                 List of renamed files.                                  |
|       type_changed_files       | string |                         List of files that had<br>type changes.                         |
|         unknown_files          | string |                                 List of unknown files.                                  |
|         unmerged_files         | string |                                 List of unmerged files.                                 |

<!-- AUTO-DOC-OUTPUT:END -->
