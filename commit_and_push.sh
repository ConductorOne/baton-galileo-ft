#!/bin/bash
cd /Users/laurenleach/go/src/github.com/ConductorOne/baton-galileo-ft

# Add new files
git add .github/workflows/capabilities_and_config.yaml
git add pkg/config/conf.gen.go
git add pkg/config/config.go
git add pkg/config/gen/gen.go

# Commit all changes
git commit -am "Add containerization support

- Move config from cmd to pkg/config package
- Add generated config struct (GalileoFt)
- Update connector to use V2 interface (ResourceSyncerV2)
- Update resource syncers to use SyncOpAttrs instead of pagination.Token
- Update baton-sdk to v0.7.10
- Update Makefile with code generation
- Update CI workflow to trigger on main branch pushes
- Add capabilities_and_config workflow
- Remove main and capabilities workflows (replaced)
- Remove lambda: false from release workflow

Co-Authored-By: Claude Sonnet 4.5 <noreply@anthropic.com>"

# Push branch
git push -u origin containerize

echo "Done!"
