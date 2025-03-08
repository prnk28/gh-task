# gh-pm

A GitHub CLI extension to find and execute tasks from Taskfiles configured in your `.github` repo.

## Quickstart

**1. Create a `.github` repo in your GitHub organization or user account.**

```bash
gh repo create <ORG/USER>/.github --public
```

**2. Add a `Taskfile.yml` to the root of your `.github` repo.**

```bash
gh repo clone <ORG/USER>/.github
cd .github
task --init
git add Taskfile.yml
git commit -m "Add Taskfile.yml"
git push
```

**3. Install the `gh-task` extension.**

```bash
gh extension install prnk28/gh-task
echo "alias ght='gh task'" >> ~/.zshrc
source ~/.zshrc
```

**4. Run your commands using the `ght` alias or `gh task`.**

```bash
gh task # or ght with alias
> # Returns "Hello World!"
```

## Usage
