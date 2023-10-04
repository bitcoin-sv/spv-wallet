# Code Standards & Contributing Guidelines

- [Code Standards \& Contributing Guidelines](#code-standards--contributing-guidelines)
  - [Most important rules - Quick Checklist](#most-important-rules---quick-checklist)
  - [1. Code style and formatting - official guidelines](#1-code-style-and-formatting---official-guidelines)
    - [1.1 In GO applications or libraries, we follow the official guidelines](#11-in-go-applications-or-libraries-we-follow-the-official-guidelines)
      - [Additional useful resources with GO recommendations, best practices and the common mistakes](#additional-useful-resources-with-go-recommendations-best-practices-and-the-common-mistakes)
  - [2. Code Rules](#2-code-rules)
    - [2.1. Self-documenting code](#21-self-documenting-code)
      - [As a Developer](#as-a-developer)
      - [As a PR Reviewer](#as-a-pr-reviewer)
    - [2.2. Tests](#22-tests)
      - [Principle](#principle)
      - [Guidelines for Writing Tests](#guidelines-for-writing-tests)
        - [Pull request title with a scope and task number](#pull-request-title-with-a-scope-and-task-number)
    - [3.3 Branching](#33-branching)
      - [Choosing branch names](#choosing-branch-names)
      - [Descriptiveness](#descriptiveness)
      - [Include Issue Number](#include-issue-number)
      - [Deleting Branches After Merging](#deleting-branches-after-merging)
      - [Remove Remote Branches](#remove-remote-branches)
      - [Recommendation: Clean Local Branches](#recommendation-clean-local-branches)
  - [4. Documentation Code Standards](#4-documentation-code-standards)
    - [4.1 Overview](#41-overview)
    - [4.2 Principles](#42-principles)
    - [4.3 Feature Documentation](#43-feature-documentation)
      - [4.3.1 Necessity](#431-necessity)
      - [4.3.2 Examples](#432-examples)
    - [4.4 External Features](#44-external-features)
    - [4.5 Markdown usage](#45-markdown-usage)
    - [4.5 Conclusion](#45-conclusion)

## Most important rules - Quick Checklist

- [ ] Follow official Go guidelines for style and formatting.
- [ ] Write self-documenting code and minimize comments.
- [ ] Ensure comprehensive test coverage including happy and error paths.
- [ ] Provide meaningful and constructive code reviews.
- [ ] Adhere to Conventional Commits for commit messages and PR naming.
- [ ] Document every feature adequately, especially for open-source projects.
- [ ] Keep documentation clear, concise, up-to-date, and accessible.
- [ ] Branching - choose consistent naming conventions, include issue number, delete branches after merging.

## 1. Code style and formatting - official guidelines

### 1.1 In GO applications or libraries, we follow the official guidelines

- [Effective Go](https://go.dev/doc/effective_go) - official Go guidelines
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments) - official Go code review comments
- [Go Examples](https://pkg.go.dev/testing#hdr-Examples) - official Go examples - used in libraries to explain how to use their exposed features
- [Go Test](https://pkg.go.dev/testing) - official Go testing package & recommendations
- [Go Linter](https://golangci-lint.run/) - golangci-lint - only codestyle checks

 > Our current linter configuration is in the `.golangci.yml` file.

#### Additional useful resources with GO recommendations, best practices and the common mistakes

- [Go Styles by Google](https://google.github.io/styleguide/go/) - Google's Go Style Guide
- [Uber Go Style Guide](https://github.com/uber-go/guide/blob/master/style.md) - Uber's Go Style Guide
- [Go Common Mistakes](https://github.com/golang/go/wiki/CommonMistakes) - Common Mistakes in Go

## 2. Code Rules

### 2.1. Self-documenting code

#### As a Developer

- Refactoring tasks should be strategically undertaken; not every piece of code warrants modification. If a segment of code, despite comments, remains untouched and unproblematic, it is   likely fulfilling its purpose effectively.
- When embarking on feature additions or bug fixes, and encountering commented code, consider extracting functions and creating a PR before proceeding with the primary task. This helps maintain code clarity and function.
- Should you come across a superfluous comment during your changes, do take the initiative to remove it and create a PR for that. This practice contributes to keeping the repository neat and well-maintained.
- Endeavor to minimize the addition of comments when making any changes to the code.

#### As a PR Reviewer

- Be vigilant of newly added comments during reviews. If a comment appears unnecessary, uninformative, or could be replaced with a function, do not hesitate to highlight this.
- Assess the meaningfulness and clarity of function names, ensuring they contribute to self-documenting code.

### 2.2. Tests

#### Principle

Developers are required to diligently cover their changes with tests and organize these tests with caution. A well-structured set of tests serves as a safety net, facilitating the swift identification and resolution of issues.

#### Guidelines for Writing Tests

1. **Readability**: Maintain test readability through the use of descriptive test names. When possible, group cases into table-driven tests, exercising discretion to avoid excessiveness. Avoid code branching in tests; instead, split differing scenarios into separate test cases.

2. **Comprehensive Coverage**: Ensure comprehensive test coverage, including error paths. Prioritize adding use cases for well-known errors that could be returned from your code.

3. **Structural Consistency**: Adopt the "Given-When-Then" structure in automatic tests to enhance readability and consistency. This structure assists team members in following the test flow and understanding its purpose, thus facilitating smoother onboarding for new members and maintaining uniformity across tests.

   - **Given**: Set the stage by preparing inputs, mocks, and application state.
   - **When**: Trigger the action or function under test.
   - **Then**: Verify if the outcomes match the expectations.

   ```go
   func TestSomethingVeryUsefulIsHappening(t *testing.T) {
       //given
       ... // prepare inputs, mocks, application state
   
       //when
       ... // call the function you're actually testing
   
       //then
       ... // check expectations about the function output
   }
    ```

4. **Test Isolation**: Ensure test isolation by avoiding the use of global variables and shared state. Each test should be independent and not rely on the execution of other tests. If a test requires a shared state, use a setup function to create the state before each test.

5. **Test Data**: Avoid random data in tests. Instead, use predefined data to ensure test consistency and reproducibility. If random data is required, use a seed to ensure the same data is generated each time the test is run.

6. **Test Cases**: If you are writing a public functions - it should be covered by tests. We should have test cases for all possible scenarios. That means that **we should have tests for all possible errors** that can be returned from the function.
Of course not only error paths should be covered - **we should highlight the happy path as well**.

7. **Testing private (unexported) functions**: When testing private functions, we should test them through the public (exported) functions that use them. The exception is when the private function is too complex to be tested through the public function. Good example is for example a function that is implementing a complex algorithm. In this case we should test the private function directly.

### 2.3 Code Review

#### Guidelines for Code Review

1. **Constructive Feedback**: Reviewers should provide constructive feedback, highlighting both strengths and areas of improvement. Comments should be clear, concise, and related to code structure, functionality, or style.

2. **Coding Standards**: Both authors and reviewers should ensure the code adheres to established coding standards, including formatting, naming conventions, and best practices as outlined in official Golang guidelines.

3. **Testing**: Reviewers should ensure that the submitted code is accompanied by adequate tests, covering both happy paths and error paths. Check the readability, descriptiveness, and structure of tests.

4. **Error Handling**: Pay special attention to error handling within the code. Ensure that errors are not ignored, logged appropriately, and presented to the user in a user-friendly format when applicable.

5. **Performance**: Review the code for any potential performance issues, such as inefficient loops, unnecessary allocations, or misuse of concurrency.

6. **Dependency Management**: Ensure that any new dependencies are necessary, appropriately versioned, and have been vetted for performance and security.

7. **Security Practices**: Ensure that the code follows secure coding practices and avoids common vulnerabilities.
   Good checklist for security practices can be found [here](https://owasp.org/www-project-top-ten/).

8. **Documentation**: Confirm that the code is well-documented, including comments, function/method descriptions, and module-level documentation as necessary.

9. **Efficiency and Readability**: The code should be efficient and readable. Reviewers should look for any code smells, overly complex functions, and ensure the use of idiomatic Go patterns.

10. **Responsiveness**: Both authors and reviewers should be timely in their responses. Authors should address all review comments, and reviewers should re-review changes promptly.

#### 2.3.3 Code Review Checklist

- [ ] Does the code adhere to the project’s coding standards?
- [ ] Are there sufficient tests, and do they cover a variety of cases?
- [ ] Is error handling comprehensive and user-friendly?
- [ ] Are there any performance concerns in the code?
- [ ] Have new dependencies been appropriately vetted?
- [ ] Does the code follow secure coding practices and avoid common vulnerabilities?
- [ ] Is the code well(self)-documented, with clear variable, function naming?
- [ ] Is the code efficient, readable, and free of code smells?
- [ ] Have all review comments been addressed in a timely manner?
- [ ] Do you understand the code and it's purpose?

This checklist serves as a guide to both authors and reviewers to ensure a thorough and effective code review process.

## 3. Contributing

### 3.1 Pull Requests && Issues

We have separate templates for Pull Requests and Issues. Please use them when creating a new PR or Issue.

### 3.2 Conventional Commits & Pull Requests Naming

#### 3.2.1 Overview

In an effort to maintain clarity and coherence in our commit history, we are adopting the Conventional Commits style for all commit messages across our repositories. This uniform format not only enhances the readability of our commit history but also facilitates automated tools in generating changelogs and extracting valuable information effectively.

#### 3.2.2 Structure

Conventional Commits follow a structured format: `type(scope): description`, where:

- `type`: Represents the nature of the commit (e.g., feat, fix, chore).
- `scope`: Denotes the relevant module or issue.
- `description`: Provides a brief explanation of the change.

When introducing breaking changes, an `!` should be appended after the `type/scope`: `feat(#123)!: introduce a breaking change`.

#### 3.2.3 Types

- `feat`: Utilized when introducing a new feature to the codebase.
- `fix`: Employed when resolving a bug or issue in the code.
- `docs`: Designated for commits involving documentation changes, such as updating README files or adding comments.
- `style`: Applied to commits focusing on code style and formatting, without altering the code's functionality.
- `refactor`: Used for code changes that neither introduce new features nor fix bugs, but improve the code structure or design.
- `test`: Assigned to commits pertaining to the addition, modification, or refactoring of tests.
- `chore`: For changes related to build processes, local development, or other maintenance tasks.
- `perf`: Employed when enhancing the performance of the codebase.
- `revert`: Marked for commits that revert a previous change.
- `ci`: Applied to changes concerning the Continuous Integration (CI) configuration or scripts.
- `deps`: Used when updating or modifying dependencies.

#### 3.2.4 Conventional Commits - Automatic Versioning

In our repositories, we use Conventional Commits to automatically generate the version number for our releases.

It works like this:

`fix: which represents bug fixes, and correlates to a SemVer patch.`
`feat: which represents a new feature, and correlates to a SemVer minor.`
`feat!:, or fix!:, refactor!:, etc., which represent a breaking change (indicated by the !) and will result in a SemVer major.`

Real life example:

`feat(#123)!: introduce breaking change - 1.0.0 -> 2.0.0`
`feat(#124): introduce new feature - 2.0.0 -> 2.1.0`
`fix(#125): fix a bug - 2.1.0 -> 2.1.1`

Given a version number MAJOR.MINOR.PATCH, increment the:

MAJOR version when you make incompatible API changes
MINOR version when you add functionality in a backward compatible manner
PATCH version when you make backward compatible bug fixes
Additional labels for pre-release and build metadata are available as extensions to the MAJOR.MINOR.PATCH format.

More about Semantic Versioning can be found [here](https://semver.org/).

#### 3.2.5 Scope

We have standardized the use of JIRA/GitHub issue numbers as the `scope` in commits within our team. This practice aids in easily tracing the origin of changes.

In the absence of an existing issue for your changes, please create one in the client’s JIRA system. If the change is not client-related, establish a GitHub issue in the repository.

#### 3.2.6 Further Reading

Additional information and guidelines on Conventional Commits can be found [here](https://www.conventionalcommits.org/en/v1.0.0/).

#### 3.2.7 Examples

##### Commit message with scope

Good example:

```bash
feat: add possibility to create a new user by admin (#BUX-123)
```

Bad example:

```bash
debugo feature - checkpoint full work
```

##### Pull request title with a scope and task number

> feat(#123): add new feature

### 3.3 Branching

#### Choosing branch names

- Choose consistent naming conventions. Common practices include:
  - `feature/feature-name`
  - `bugfix/issue-or-bug-name`
  - `hotfix/hotfix-name`
  - `chore/task-name`
  - `refactor/refactor-name`

#### Descriptiveness

- Branch names should be descriptive and represent the task/feature at hand.
- Use hyphens to separate words for readability, e.g., `feature/add-login-button`.

#### Include Issue Number

- If applicable, include the issue number in the branch name for easy tracking, e.g., `feature/123-add-login-button`.

#### Deleting Branches After Merging

#### Remove Remote Branches

- Once a PR has been merged, delete the remote branch to keep the repository clean.
- GitHub provides a button to delete the branch once the PR is merged.

#### Recommendation: Clean Local Branches

- Regularly prune local branches that have been deleted remotely with `git fetch -p && git branch -vv | grep 'origin/.*: gone]' | awk '{print $1}' | xargs git branch -d`.

## 4. Documentation Code Standards

### 4.1 Overview

A well-documented codebase is pivotal for both internal development and external contributions, especially for open-source projects that expose functionalities for public use. Comprehensive documentation, supplemented with examples where necessary, ensures that every feature is easily understandable, usable, and maintainable.

### 4.2 Principles

- **Clarity and Conciseness**: Documentation should be clear, concise, and focused, providing necessary information without unnecessary complexity or verbosity.
- **Accessibility**: Documentation should be accessible to developers of varying skill levels, enabling both novices and experts to understand the codebase and its features.
- **Up-to-Date**: Documentation must be kept current, reflecting the latest changes and developments in the codebase to avoid misinformation and confusion.

### 4.3 Feature Documentation

#### 4.3.1 Necessity

Every feature developed should be accompanied by adequate documentation. The necessity for documentation becomes even more pronounced for open-source projects, where clear instructions and examples facilitate easier adoption and contribution from the community.

#### 4.3.2 Examples

- **Inclusion of Examples**: Where applicable, documentation should include practical examples demonstrating the feature’s usage and benefits. Examples act as a practical guide, aiding developers in understanding and implementing the feature correctly.
- **Clarity of Examples**: Examples should be clear, concise, and relevant, illustrating the functionality of the feature effectively.

### 4.4 External Features

For projects exposing external features:

- **Comprehensive Guides**: Ensure the creation of comprehensive guides detailing the utilization of exposed features, their benefits, and any potential configurations or customizations.
- **Community Engagement**: Encourage community members to contribute to documentation by providing feedback, suggestions, and improvements. This collaborative approach enriches the documentation quality and breadth.

### 4.5 Markdown usage

We should write documentation in Markdown format. It allows us to write documentation in a simple and readable way. It's also easy to convert Markdown to HTML or PDF or create a website from it.

[Markdown Guide](markdownguide.org) - Comprehensive guide to Markdown syntax.

### 4.5 Conclusion

Adhering to documentation code standards is integral for maintaining a healthy, understandable, and contributable codebase. By ensuring every feature is well-documented, with the inclusion of clear examples where necessary, we foster a conducive environment for development and community engagement, particularly in open-source projects.
