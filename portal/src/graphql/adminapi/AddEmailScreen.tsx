import React, { useCallback, useMemo, useState } from "react";
import { useParams } from "react-router-dom";
import deepEqual from "deep-equal";
import { FormattedMessage } from "@oursky/react-messageformat";

import {
  ModifiedIndicatorPortal,
  ModifiedIndicatorWrapper,
} from "../../ModifiedIndicatorPortal";
import NavBreadcrumb from "../../NavBreadcrumb";
import ShowError from "../../ShowError";
import FormTextField from "../../FormTextField";
import ShowLoading from "../../ShowLoading";
import AddIdentityForm from "./AddIdentityForm";
import { passwordFieldErrorRules } from "../../PasswordField";
import ShowUnhandledValidationErrorCause from "../../error/ShowUnhandledValidationErrorCauses";
import { useCreateLoginIDIdentityMutation } from "./mutations/createIdentityMutation";
import { useTextField } from "../../hook/useInput";
import { PortalAPIAppConfig } from "../../types";
import { UserQuery_node_User } from "./query/__generated__/UserQuery";
import { useUserQuery } from "./query/userQuery";
import { useAppConfigQuery } from "../portal/query/appConfigQuery";
import { FormContext } from "../../error/FormContext";
import { useValidationError } from "../../error/useValidationError";
import { useGenericError } from "../../error/useGenericError";

import styles from "./AddEmailScreen.module.scss";

interface AddEmailFormData {
  email: string;
  password: string;
}

interface AddEmailFormProps {
  appConfig: PortalAPIAppConfig | null;
  user: UserQuery_node_User | null;
}

const AddEmailForm: React.FC<AddEmailFormProps> = function AddEmailForm(
  props: AddEmailFormProps
) {
  const { appConfig, user } = props;
  const { userID } = useParams();

  const {
    createIdentity,
    loading: creatingIdentity,
    error: createIdentityError,
  } = useCreateLoginIDIdentityMutation(userID);

  const initialFormData = useMemo(() => {
    return {
      email: "",
      password: "",
    };
  }, []);
  const [formData, setFormData] = useState<AddEmailFormData>(initialFormData);
  const { email, password } = formData;

  const [localValidationErrorMessage, setLocalViolationErrorMessage] = useState<
    string | undefined
  >(undefined);

  const { onChange: onEmailChange } = useTextField((value) => {
    setFormData((prev) => ({ ...prev, email: value }));
  });
  const { onChange: onPasswordChange } = useTextField((value) => {
    setFormData((prev) => ({ ...prev, password: value }));
  });

  const isFormModified = useMemo(() => {
    return !deepEqual(initialFormData, formData);
  }, [initialFormData, formData]);

  const resetForm = useCallback(() => {
    setFormData(initialFormData);
    setLocalViolationErrorMessage(undefined);
  }, [initialFormData]);

  const {
    unhandledCauses: rawUnhandledCauses,
    otherError,
    value: formContextValue,
  } = useValidationError(createIdentityError);

  const {
    errorMessageMap,
    unrecognizedError,
    unhandledCauses,
  } = useGenericError(otherError, rawUnhandledCauses, [
    {
      reason: "InvariantViolated",
      kind: "DuplicatedIdentity",
      errorMessageID: "AddEmailScreen.error.duplicated-email",
      field: "email",
    },
    ...passwordFieldErrorRules,
  ]);

  return (
    <FormContext.Provider value={formContextValue}>
      <ModifiedIndicatorPortal
        resetForm={resetForm}
        isModified={isFormModified}
      />
      {unrecognizedError && (
        <div className={styles.error}>
          <ShowError error={unrecognizedError} />
        </div>
      )}
      <ShowUnhandledValidationErrorCause causes={unhandledCauses} />
      <AddIdentityForm
        className={styles.content}
        appConfig={appConfig}
        user={user}
        password={password}
        onPasswordChange={onPasswordChange}
        passwordFieldErrorMessage={
          localValidationErrorMessage ?? errorMessageMap.password
        }
        loginIdKey="email"
        loginId={email}
        isFormModified={isFormModified}
        createIdentity={createIdentity}
        creatingIdentity={creatingIdentity}
        onLocalErrorMessageChange={setLocalViolationErrorMessage}
        loginIdField={
          <FormTextField
            jsonPointer=""
            parentJSONPointer=""
            fieldName="email"
            fieldNameMessageID="AddEmailScreen.email.label"
            className={styles.emailField}
            value={email}
            onChange={onEmailChange}
            errorMessage={errorMessageMap.email}
          />
        }
      />
    </FormContext.Provider>
  );
};

const AddEmailScreen: React.FC = function AddEmailScreen() {
  const { userID, appID } = useParams();

  const {
    user,
    loading: loadingUser,
    error: userError,
    refetch: refetchUser,
  } = useUserQuery(userID);
  const {
    effectiveAppConfig,
    loading: loadingAppConfig,
    error: appConfigError,
    refetch: refetchAppConfig,
  } = useAppConfigQuery(appID);

  const navBreadcrumbItems = useMemo(() => {
    return [
      { to: "../../..", label: <FormattedMessage id="UsersScreen.title" /> },
      { to: "../", label: <FormattedMessage id="UserDetailsScreen.title" /> },
      { to: ".", label: <FormattedMessage id="AddEmailScreen.title" /> },
    ];
  }, []);

  if (loadingUser || loadingAppConfig) {
    return <ShowLoading />;
  }

  if (userError != null) {
    return <ShowError error={userError} onRetry={refetchUser} />;
  }

  if (appConfigError != null) {
    return <ShowError error={appConfigError} onRetry={refetchAppConfig} />;
  }

  return (
    <div className={styles.root}>
      <ModifiedIndicatorWrapper>
        <NavBreadcrumb
          className={styles.breadcrumb}
          items={navBreadcrumbItems}
        />
        <AddEmailForm appConfig={effectiveAppConfig} user={user} />
      </ModifiedIndicatorWrapper>
    </div>
  );
};

export default AddEmailScreen;
