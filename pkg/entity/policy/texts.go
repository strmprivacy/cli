package policy

import "strmprivacy/strm/pkg/util"

var longCreateDoc = util.LongDocsUsage(`Create a Policy

A policy has the following attributes
* name: the name of a policy. This must be unique within one organization.
* description: a description of the policy; what sort of data pipelines
  would be subject to this policy?
* retention: the number of days that encryption keys created under this
  policy should be kept. This might be a minimum or a maximum...
* legal grounds: a legal text or ruling that identifies why the organization
  created this policy
* state: draft, active or archived. Policies can only be used in pipelines
  when they're in active state. They can still be modified while in draft.
  Deletion of policies is not allowed for active policies.
`)

var policyExample = util.DedentTrim(`
			strm get policy "1 year" or strm get policy --id=34c4709e-b8bc-4b45-aa5a-883f471869e3
				Name: 1 year
				Id: 5c8e653a-8102-4444-ac15-a3d1aa0ff109
				Description:
				Retention(days): 365
				Legal Grounds:
				State: STATE_ACTIVE

			strm get policy --get-default-policy
				Name: 7 years
				Id:
				Description: Default 7 year retention
				Retention(days): 2556
				Legal Grounds:
				State: STATE_ACTIVE

		`)

var listExample = util.DedentTrim(`
strm list policies
 NAME            DESCRIPTION    RETENTION(DAYS)   LEGAL GROUNDS           STATE

 1 year                                     365                    STATE_DRAFT
 2 long years    2 whole years              730   GDPR             STATE_ACTIVE
`)
