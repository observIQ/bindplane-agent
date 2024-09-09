package githubreceiver

type gitHubLog struct {
	user *gitHubUserLog
    org *gitHubOrganizationLog
    enterprise *gitHubEnterpriseLog
}

type ActorLocation struct {
	CountryCode string `json:"country_code"`
}

type gitHubEnterpriseLog struct {
    Timestamp                   int64         `json:"@timestamp"`
    DocumentID                  string        `json:"document_id"`
    Action                      string        `json:"action"`
    Actor                       string        `json:"actor"`
    ActorID                     int64         `json:"actor_id"`
    ActorIsBot                  bool          `json:"actor_is_bot,omitempty"`
    ActorIP                     string        `json:"actor_ip,omitempty"`
    ActorLocation               ActorLocation `json:"actor_location,omitempty"`
    CountryCode                 string        `json:"country_code,omitempty"`
    Business                    string        `json:"business,omitempty"`
    BusinessID                  int64         `json:"business_id,omitempty"`
    CreatedAt                   int64         `json:"created_at"`
    OperationType               string        `json:"operation_type,omitempty"`
	UserAgent				   string        `json:"user_agent,omitempty"`
    ActorLogin                  string        `json:"actor_login,omitempty"`
    ActorLocationCountryCode    string        `json:"actor_location_country_code,omitempty"`
    ActorAvatarURL              string        `json:"actor_avatar_url,omitempty"`
    ActorIsEnterpriseOwner      bool          `json:"actor_is_enterprise_owner,omitempty"`
    ActorUserAgent              string        `json:"actor_user_agent,omitempty"`
    Repository                  string        `json:"repository,omitempty"`
    RepoPrivate                 bool          `json:"repo_private,omitempty"`
    RepoVisibility              string        `json:"repo_visibility,omitempty"`
    TargetLogin                 string        `json:"target_login,omitempty"`
    TargetType                  string        `json:"target_type,omitempty"`
    Team                        string        `json:"team,omitempty"`
    TeamID                      int64         `json:"team_id,omitempty"`
    TeamSlug                    string        `json:"team_slug,omitempty"`
    Enterprise                  string        `json:"enterprise,omitempty"`
    EnterpriseID                int64         `json:"enterprise_id,omitempty"`
    EnterpriseSlug              string        `json:"enterprise_slug,omitempty"`
    User                        string        `json:"user,omitempty"`
    UserID                      int64         `json:"user_id,omitempty"`
    UserLogin                   string        `json:"user_login,omitempty"`
    Permission                  string        `json:"permission,omitempty"`
    Ref                         string        `json:"ref,omitempty"`
    Branch                      string        `json:"branch,omitempty"`
    Environment                 string        `json:"environment,omitempty"`
    Workflow                    string        `json:"workflow,omitempty"`
    Deployment                  string        `json:"deployment,omitempty"`
    RunID                       int64         `json:"run_id,omitempty"`
    InstallationID              int64         `json:"installation_id,omitempty"`
    InvitationID                int64         `json:"invitation_id,omitempty"`
    Integration                 string        `json:"integration,omitempty"`
    IntegrationID               int64         `json:"integration_id,omitempty"`
    ExternalURL                 string        `json:"external_url,omitempty"`
    DocumentationURL            string        `json:"documentation_url,omitempty"`
    EnvironmentName             string        `json:"environment_name,omitempty"`
    JobName                     string        `json:"job_name,omitempty"`
    JobStatus                   string        `json:"job_status,omitempty"`
    OrganizationUpgrade         bool          `json:"organization_upgrade,omitempty"`
    Plan                        string        `json:"plan,omitempty"`
    BillingEmail                string        `json:"billing_email,omitempty"`
    AuditLogStreamSink          string        `json:"audit_log_stream_sink,omitempty"`
    AuditLogStreamResult        string        `json:"audit_log_stream_result,omitempty"`
    DeploymentEnvironment       string        `json:"deployment_environment,omitempty"`
    Member                      string        `json:"member,omitempty"`
    MemberLogin                 string        `json:"member_login,omitempty"`
    SSHKey                      string        `json:"ssh_key,omitempty"`
    SSHKeyID                    int64         `json:"ssh_key_id,omitempty"`
    TargetID                    int64         `json:"target_id,omitempty"`
    RepositoryID                int64         `json:"repository_id,omitempty"`
    RepositoryPublic            bool          `json:"repository_public,omitempty"`
    RepoOwner                   string        `json:"repo_owner,omitempty"`
    RepoOwnerID                 int64         `json:"repo_owner_id,omitempty"`
    ProtectedBranch             bool          `json:"protected_branch,omitempty"`
    RefName                     string        `json:"ref_name,omitempty"`
    OrganizationBillingEmail    string        `json:"organization_billing_email,omitempty"`
    PreviousPermission          string        `json:"previous_permission,omitempty"`
    HookID                      int64         `json:"hook_id,omitempty"`
    HookURL                     string        `json:"hook_url,omitempty"`
    HookName                    string        `json:"hook_name,omitempty"`
    BranchProtection            string        `json:"branch_protection,omitempty"`
    WorkflowRunID               int64         `json:"workflow_run_id,omitempty"`
    WorkflowFileName            string        `json:"workflow_file_name,omitempty"`
    WorkflowFilePath            string        `json:"workflow_file_path,omitempty"`
    RunAttempt                  int64         `json:"run_attempt,omitempty"`
    WorkflowRunStartedAt        int64         `json:"workflow_run_started_at,omitempty"`
    WorkflowRunConclusion       string        `json:"workflow_run_conclusion,omitempty"`
    SSHKeyTitle                 string        `json:"ssh_key_title,omitempty"`
    SSHKeyFingerprint           string        `json:"ssh_key_fingerprint,omitempty"`
    OAuthTokenID                int64         `json:"oauth_token_id,omitempty"`
    OAuthTokenName              string        `json:"oauth_token_name,omitempty"`
    ApplicationID               int64         `json:"application_id,omitempty"`
    ApplicationName             string        `json:"application_name,omitempty"`
    License                     string        `json:"license,omitempty"`
    LicenseExpiry               int64         `json:"license_expiry,omitempty"`
    SAMLNameID                  string        `json:"saml_name_id,omitempty"`
    SAMLNameIDEmail             string        `json:"saml_name_id_email,omitempty"`
    SAMLNameIDUser              string        `json:"saml_name_id_user,omitempty"`
    SSOID                       int64         `json:"sso_id,omitempty"`
    SSOName                     string        `json:"sso_name,omitempty"`
    TwoFAEnforcement            bool          `json:"2fa_enforcement,omitempty"`
    TwoFAType                   string        `json:"2fa_type,omitempty"`
    OrgName                     string        `json:"org_name,omitempty"`
    OrgRole                     string        `json:"org_role,omitempty"`
    OrgLogin                    string        `json:"org_login,omitempty"`
    ActionDescription           string        `json:"action_description,omitempty"`
    LDAPDN                      string        `json:"ldap_dn,omitempty"`
    MFA                         string        `json:"mfa,omitempty"`
    MFAEnrollment               bool          `json:"mfa_enrollment,omitempty"`
    Name                        string        `json:"name,omitempty"`
    Org                         string        `json:"org,omitempty"`
    OrgID                       int64         `json:"org_id,omitempty"`
    OwnerType                   string        `json:"owner_type,omitempty"`
}


type gitHubOrganizationLog struct {
    Timestamp                   int64         `json:"@timestamp"`
    Action                      string        `json:"action"`
    Actor                       string        `json:"actor"`
    ActorID                     int64         `json:"actor_id"`
    ActorIP                     string        `json:"actor_ip,omitempty"`
    ActorLocation               ActorLocation `json:"actor_location,omitempty"`
    CountryCode                 string        `json:"country_code,omitempty"`
    ActorLogin                  string        `json:"actor_login,omitempty"`
    ActorLocationCountryCode    string        `json:"actor_location_country_code,omitempty"`
    ActorAvatarURL              string        `json:"actor_avatar_url,omitempty"`
    Business                    string        `json:"business,omitempty"`
    BusinessID                  int64         `json:"business_id,omitempty"`
    CreatedAt                   int64         `json:"created_at"`
    OperationType               string        `json:"operation_type,omitempty"`
    Repository                  string        `json:"repository,omitempty"`
    RepoPrivate                 bool          `json:"repo_private,omitempty"`
    RepoVisibility              string        `json:"repo_visibility,omitempty"`
    TargetLogin                 string        `json:"target_login,omitempty"`
    TargetType                  string        `json:"target_type,omitempty"`
    Team                        string        `json:"team,omitempty"`
    TeamID                      int64         `json:"team_id,omitempty"`
    TeamName                    string        `json:"team_name,omitempty"`
    TeamSlug                    string        `json:"team_slug,omitempty"`
    Org                         string        `json:"org,omitempty"`
    OrgID                       int64         `json:"org_id,omitempty"`
    OrgLogin                    string        `json:"org_login,omitempty"`
    User                        string        `json:"user,omitempty"`
    UserID                      int64         `json:"user_id,omitempty"`
    UserLogin                   string        `json:"user_login,omitempty"`
    Permission                  string        `json:"permission,omitempty"`
    Ref                         string        `json:"ref,omitempty"`
    Branch                      string        `json:"branch,omitempty"`
    Environment                 string        `json:"environment,omitempty"`
    Workflow                    string        `json:"workflow,omitempty"`
    Deployment                  string        `json:"deployment,omitempty"`
    RunID                       int64         `json:"run_id,omitempty"`
    InstallationID              int64         `json:"installation_id,omitempty"`
    InvitationID                int64         `json:"invitation_id,omitempty"`
    Integration                 string        `json:"integration,omitempty"`
    IntegrationID               int64         `json:"integration_id,omitempty"`
    ExternalURL                 string        `json:"external_url,omitempty"`
    DocumentationURL            string        `json:"documentation_url,omitempty"`
    EnvironmentName             string        `json:"environment_name,omitempty"`
    JobName                     string        `json:"job_name,omitempty"`
    JobStatus                   string        `json:"job_status,omitempty"`
    OrganizationUpgrade         bool          `json:"organization_upgrade,omitempty"`
    Plan                        string        `json:"plan,omitempty"`
    BillingEmail                string        `json:"billing_email,omitempty"`
    AuditLogStreamSink          string        `json:"audit_log_stream_sink,omitempty"`
    AuditLogStreamResult        string        `json:"audit_log_stream_result,omitempty"`
    DeploymentEnvironment       string        `json:"deployment_environment,omitempty"`
    Enterprise                  string        `json:"enterprise,omitempty"`
    EnterpriseID                int64         `json:"enterprise_id,omitempty"`
    Member                      string        `json:"member,omitempty"`
    MemberLogin                 string        `json:"member_login,omitempty"`
    ActorIsBot                  bool          `json:"actor_is_bot,omitempty"`
    ActorEnterpriseOwner        bool          `json:"actor_enterprise_owner,omitempty"`
    SSHKey                      string        `json:"ssh_key,omitempty"`
    SSHKeyID                    int64         `json:"ssh_key_id,omitempty"`
    TargetID                    int64         `json:"target_id,omitempty"`
    RepositoryID                int64         `json:"repository_id,omitempty"`
    RepositoryPublic            bool          `json:"repository_public,omitempty"`
    RepoOwner                   string        `json:"repo_owner,omitempty"`
    RepoOwnerID                 int64         `json:"repo_owner_id,omitempty"`
    ProtectedBranch             bool          `json:"protected_branch,omitempty"`
    RefName                     string        `json:"ref_name,omitempty"`
    OrganizationBillingEmail    string        `json:"organization_billing_email,omitempty"`
    PreviousPermission          string        `json:"previous_permission,omitempty"`
    HookID                      int64         `json:"hook_id,omitempty"`
    HookURL                     string        `json:"hook_url,omitempty"`
    HookName                    string        `json:"hook_name,omitempty"`
    BranchProtection            string        `json:"branch_protection,omitempty"`
    WorkflowRunID               int64         `json:"workflow_run_id,omitempty"`
    WorkflowFileName            string        `json:"workflow_file_name,omitempty"`
    WorkflowFilePath            string        `json:"workflow_file_path,omitempty"`
    RunAttempt                  int64         `json:"run_attempt,omitempty"`
    WorkflowRunStartedAt        int64         `json:"workflow_run_started_at,omitempty"`
    WorkflowRunConclusion       string        `json:"workflow_run_conclusion,omitempty"`
    SSHKeyTitle                 string        `json:"ssh_key_title,omitempty"`
    SSHKeyFingerprint           string        `json:"ssh_key_fingerprint,omitempty"`
}

type gitHubUserLog struct {
    ID        int64       `json:"id"`
    Type      string      `json:"type"`
    Actor     Actor       `json:"actor"`
    Repo      Repository  `json:"repo"`
    Payload   interface{} `json:"payload,omitempty"`
    Org       Organization `json:"org,omitempty"`
    Public    bool        `json:"public"`
    CreatedAt int64       `json:"created_at"`

}

type Actor struct {
	ID          int64  `json:"id"`
	Login       string `json:"login"`
	DisplayLogin string `json:"display_login,omitempty"`
	GravatarID  string `json:"gravatar_id,omitempty"`
	URL         string `json:"url"`
	AvatarURL   string `json:"avatar_url"`
} 

type Organization struct {
    ID         int64  `json:"id,omitempty"`
    Login      string `json:"login,omitempty"`
    GravatarID string `json:"gravatar_id,omitempty"`
    URL        string `json:"url,omitempty"`
    AvatarURL  string `json:"avatar_url,omitempty"`
}

type Repository struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	URL  string `json:"url"`
}

type Org struct {
	ID         int64  `json:"id,omitempty"`
	Login      string `json:"login,omitempty"`
	GravatarID string `json:"gravatar_id,omitempty"`
	URL        string `json:"url,omitempty"`
	AvatarURL  string `json:"avatar_url,omitempty"`
}

type PushEvent struct {
    Ref        string `json:"ref"`
    Before     string `json:"before"`
    After      string `json:"after"`
    Commits    []Commit `json:"commits"`
    Size       int    `json:"size"`
    Pusher     User   `json:"pusher"`
    Repository Repository `json:"repository"`
}

type Commit struct {
    SHA     string `json:"sha"`
    Author  Author `json:"author"`
    Message string `json:"message"`
    Distinct bool  `json:"distinct"`
    URL     string `json:"url"`
}

type Author struct {
    Name  string `json:"name"`
    Email string `json:"email"`
}

type User struct {
    Name  string `json:"name"`
    Email string `json:"email"`
}

type PullRequestEvent struct {
    Action      string `json:"action"`
    Number      int    `json:"number"`
	PullRequest PullRequest `json:"PullRequestEvent"`
}

type PullRequest struct {
	ID                int64  `json:"id,omitempty"`
	URL               string `json:"url,omitempty"`
	HTMLURL           string `json:"html_url,omitempty"`
	DiffURL           string `json:"diff_url,omitempty"`
	PatchURL          string `json:"patch_url,omitempty"`
	IssueURL          string `json:"issue_url,omitempty"`
	CommitsURL        string `json:"commits_url,omitempty"`
	ReviewCommentsURL string `json:"review_comments_url,omitempty"`
	ReviewCommentURL  string `json:"review_comment_url,omitempty"`
	CommentsURL       string `json:"comments_url,omitempty"`
	StatusesURL       string `json:"statuses_url,omitempty"`
	Title             string `json:"title,omitempty"`
	User              struct {
		Login string `json:"login,omitempty"`
	} `json:"user,omitempty"`
	Body               string   `json:"body,omitempty"`
	CreatedAt          int64    `json:"created_at,omitempty"`
	UpdatedAt          int64    `json:"updated_at,omitempty"`
	ClosedAt           int64    `json:"closed_at,omitempty"`
	MergedAt           int64    `json:"merged_at,omitempty"`
	MergeCommitSHA     string   `json:"merge_commit_sha,omitempty"`
	Assignees          []string `json:"assignees,omitempty"`
	RequestedReviewers []string `json:"requested_reviewers,omitempty"`
}

type DeleteEvent struct {
    Ref        string `json:"ref,omitempty"`
    RefType    string `json:"ref_type,omitempty"`
    PusherType string `json:"pusher_type,omitempty"`
}

type WatchEvent struct {
    Action string `json:"action,omitempty"`
}

type ReleaseEvent struct {
    Action  string `json:"action,omitempty"`
    Release struct {
        ID              int64  `json:"id,omitempty"`
        TagName         string `json:"tag_name,omitempty"`
        TargetCommitish string `json:"target_commitish,omitempty"`
        Name            string `json:"name,omitempty"`
        Body            string `json:"body,omitempty"`
    } `json:"release,omitempty"`
}

