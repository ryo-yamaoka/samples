spanner-create:
	gcloud --project $$SPANNER_PROJECT_ID spanner instances create $$SPANNER_INSTANCE_ID --processing-units 100 --config regional-asia-northeast1 --description $$SPANNER_INSTANCE_ID
	gcloud --project $$SPANNER_PROJECT_ID spanner databases create $$SPANNER_DATABASE_ID --instance $$SPANNER_INSTANCE_ID --ddl-file schema.sql

spanner-reset:
	-gcloud --project $$SPANNER_PROJECT_ID spanner databases delete $$SPANNER_DATABASE_ID --instance $$SPANNER_INSTANCE_ID -q
	-gcloud --project $$SPANNER_PROJECT_ID spanner databases create $$SPANNER_DATABASE_ID --instance $$SPANNER_INSTANCE_ID --ddl-file schema.sql

spanner-delete:
	-gcloud --project $$SPANNER_PROJECT_ID spanner instances delete $$SPANNER_INSTANCE_ID -q
