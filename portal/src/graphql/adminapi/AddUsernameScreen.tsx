import React, { useCallback, useMemo, useState } from "react";
import { useParams } from "react-router-dom";
import { FormattedMessage } from "@oursky/react-messageformat";
import deepEqual from "deep-equal";

import { useAppConfigQuery } from "../portal/query/appConfigQuery";
import { useUserQuery } from "./query/userQuery";
import { UserQuery_node_User } from "./query/__generated__/UserQuery";
import NavBreadcrumb from "../../NavBreadcrumb";
import { passwordFieldErrorRules } from "../../PasswordField";
import AddIdentityForm from "./AddIdentityForm";
import ShowUnhandledValidationErrorCause from "../../error/ShowUnhandledValidationErrorCauses";
import FormTextField from "../../FormTextField";
import ShowLoading from "../../ShowLoading";
import ShowError from "../../ShowError";
import {
  ModifiedIndicatorPortal,
  ModifiedIndicatorWrapper,
} from "../../ModifiedIndicatorPortal";
import { useCreateLoginIDIdentityMutation } from "./mutations/createIdentityMutation";
import { useTextField } from "../../hook/useInput";
import { useValidationError } from "../../error/useValidationError";
import { FormContext } from "../../error/FormContext";
import { useGenericError } from "../../error/useGenericError";
import { PortalAPIAppConfig } from "../../types";

import styles from "./AddUsernameScreen.module.scss";

interface AddUsernameFormProps {
  appConfig: PortalAPIAppConfig | null;
  user: UserQuery_node_User | null;
}

const AddUsernameForm: React.FC<AddUsernameFormProps> = function AddUsernameForm(
  props: AddUsernameFormProps
) {
  const { appConfig, user } = props;

  const { userID } = useParams();

  const {
    createIdentity,
    loading: creatingIdentity,
    error: createIdentityError,
  } = useCreateLoginIDIdentityMutation(userID);

  const [localValidationErrorMessage, setLocalViolationErrorMessage] = useState<
    string | undefined
  >(undefined);

  const initialFormData = useMemo(() => {
    return {
      password: "",
      username: "",
    };
  }, []);

  const [formData, setFormData] = useState(initialFormData);
  const { username, password } = formData;

  const { onChange: onUsernameChange } = useTextField((value) => {
    setFormData((prev) => ({ ...prev, username: value }));
  });
  const { onChange: onPasswordChange } = useTextField((value) => {
    setFormData((prev) => ({ ...prev, password: value }));
  });

  const isFormModified = useMemo(() => {
    return !deepEqual(formData, initialFormData);
  }, [formData, initialFormData]);

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
      errorMessageID: "AddUsernameScreen.error.duplicated-username",
      field: "username",
    },
    ...passwordFieldErrorRules,
  ]);

  return (
    <FormContext.Provider value={formContextValue}>
      <ModifiedIndicatorPortal
        resetForm={resetForm}
        isModified={isFormModified}
      />
      {unrecognizedError && <ShowError error={unrecognizedError} />}
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
        loginIdKey="username"
        loginId={username}
        isFormModified={isFormModified}
        createIdentity={createIdentity}
        creatingIdentity={creatingIdentity}
        onLocalErrorMessageChange={setLocalViolationErrorMessage}
        loginIdField={
          <FormTextField
            jsonPointer=""
            parentJSONPointer=""
            fieldName="username"
            fieldNameMessageID="AddUsernameScreen.username.label"
            className={styles.usernameField}
            value={username}
            onChange={onUsernameChange}
            errorMessage={errorMessageMap.username}
          />
        }
      />
    </FormContext.Provider>
  );
};

const AddUsernameScreen: React.FC = function AddUsernameScreen() {
  const { appID, userID } = useParams();
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
      { to: ".", label: <FormattedMessage id="AddUsernameScreen.title" /> },
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
      <ModifiedIndicatorWrapper className={styles.wrapper}>
        <NavBreadcrumb items={navBreadcrumbItems} />
        <AddUsernameForm appConfig={effectiveAppConfig} user={user} />
      </ModifiedIndicatorWrapper>
    </div>
  );
};

export default AddUsernameScreen;
